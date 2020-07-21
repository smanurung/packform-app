package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HTTPResponse struct {
	CurrentPage int `json:"current_page,omitempty"`
	Total       int `json:"total,omitempty"`
	From        int `json:"from,omitempty"`
	To          int `json:"to,omitempty"`
	// NextPageURL string  `json:"next_page_url,omitempty"`
	// PrevPageURL string  `json:"prev_page_url,omitempty"`
	PerPage  int     `json:"per_page,omitempty"`
	LastPage int     `json:"last_page,omitempty"`
	Data     []Order `json:"data,omitempty"`
}

type Order struct {
	ID                 string  `db:"id"`
	OrderName          string  `json:"order_name,omitempty" db:"order_name"`
	CustomerID         string  `db:"customer_id"`
	CustomerCompany    string  `json:"customer_company,omitempty"`
	CustomerName       string  `json:"customer_name,omitempty"`
	OrderDate          string  `json:"order_date,omitempty" db:"created_at"`
	DeliveredAmount    float64 `json:"delivered_amount_float,omitempty" db:"delivered_amount"`
	DeliveredAmountStr string  `json:"delivered_amount,omitempty"`
	TotalAmount        float64 `json:"total_amount_float,omitempty" db:"amount"`
	TotalAmountStr     string  `json:"total_amount,omitempty"`
}

type Customer struct {
	UserID      string `bson:"user_id"`
	Login       string `bson:"login"`
	Password    string `bson:"password"`
	Name        string `bson:"name"`
	CompanyID   int    `bson:"company_id"`
	CreditCards string `bson:"credit_cards"`
}

type Company struct {
	CompanyID   int    `bson:"company_id"`
	CompanyName string `bson:"company_name"`
}

const (
	layoutFrom = "2006-01-02T15:04:05Z"
	layoutTo   = "Jan 2, 3:04 PM"
)

func generateHandler(db *sqlx.DB, mongodb *mongo.Database) func(w http.ResponseWriter, r *http.Request) {
	loc, err := time.LoadLocation("Australia/Brisbane")
	if err != nil {
		log.Errorln(err)
	}

	return (func(w http.ResponseWriter, r *http.Request) {

		page := r.FormValue("page")
		perPage := r.FormValue("per_page")

		offset := 0
		pageInt := 0
		perPageInt := 10
		var err error
		if page != "" && perPage != "" {
			pageInt, err = strconv.Atoi(page)
			if err != nil {
				log.Errorln(err)
			}
			perPageInt, err = strconv.Atoi(perPage)
			if err != nil {
				log.Errorln(err)
			}
			offset = (pageInt - 1) * perPageInt
		}
		log.Infoln(page, perPage, offset)

		filter := r.FormValue("filter")
		startDate := r.FormValue("start_date")
		endDate := r.FormValue("end_date")

		// where := ""
		// if filter := r.FormValue("filter"); filter != "" {
		// 	where = "where order_name ilike '%" + filter + "%'"
		// }

		var filters []string
		var args []interface{}
		idx := 1 // since $1 and $2 are for paging (limit & offset).
		if filter != "" {
			filters = append(filters, fmt.Sprintf("order_name ilike $%d", idx))
			args = append(args, "%"+filter+"%")
			idx++
		}
		if startDate != "" {
			filters = append(filters, fmt.Sprintf("DATE(created_at) >= $%d", idx))
			args = append(args, startDate)
			idx++
		}
		if endDate != "" {
			filters = append(filters, fmt.Sprintf("DATE(created_at) <= $%d", idx))
			args = append(args, endDate)
			idx++
		}

		where := ""
		if len(filters) > 0 {
			where = "where " + strings.Join(filters, " and ")
		}

		query := `
			with some_orders as (
				select id, order_name, created_at, customer_id
				from orders
				` + where + `	order by created_at desc
				` + fmt.Sprintf("limit $%d offset $%d", idx, idx+1) +
			`),
			some_order_items as (
				select A.id, order_name, created_at, customer_id, B.id as order_item_id, price_per_unit, quantity
				from some_orders as A left join order_items as B
				on A.id = B.order_id
			)
			select order_name, created_at, customer_id, sum(price_per_unit*quantity) as amount, sum(coalesce(price_per_unit*delivered_quantity,0)) as delivered_amount
			from some_order_items as A left join delivery as B
			on A.order_item_id = B.order_item_id
			group by order_name, created_at, customer_id
			order by created_at desc
		`

		log.Infoln(query)
		log.Infoln(append(args, perPage, offset))

		var orders []Order
		err = db.Select(&orders, query, append(args, perPage, offset)...)
		if err != nil {
			log.Errorln(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// count query
		query = "select count(1) from orders " + where
		log.Infoln(query)

		var total int
		err = db.Get(&total, query, args...)
		if err != nil {
			log.Errorln(err)
		}
		lastPage := total / perPageInt

		customerColl := mongodb.Collection("customers")
		companiesColl := mongodb.Collection("customer_companies")

		var data []Order
		for _, o := range orders {
			log.Infoln(o)

			var customer Customer
			filterCustomer := bson.D{{"user_id", o.CustomerID}}
			err = customerColl.FindOne(context.TODO(), filterCustomer).Decode(&customer)
			if err != nil {
				log.Errorln(err)
			}

			var company Company
			filterCompany := bson.D{{"company_id", customer.CompanyID}}
			err = companiesColl.FindOne(context.TODO(), filterCompany).Decode(&company)
			if err != nil {
				log.Errorln(err)
			}

			parsedTime, err := time.Parse(layoutFrom, o.OrderDate)
			if err != nil {
				log.Errorln(err)
			}

			o.CustomerCompany = company.CompanyName
			o.CustomerName = customer.Name
			o.OrderDate = parsedTime.In(loc).Format(layoutTo)
			o.TotalAmountStr = fmt.Sprintf("$%.2f", o.TotalAmount)

			o.DeliveredAmountStr = "-"
			if o.DeliveredAmount > 0 {
				o.DeliveredAmountStr = fmt.Sprintf("$%.2f", o.DeliveredAmount)
			}

			data = append(data, o)
		}

		resp := HTTPResponse{
			CurrentPage: pageInt,
			Total:       total,
			From:        offset + 1,
			To:          offset + perPageInt,
			// NextPageURL: "",
			// PrevPageURL: "",
			PerPage:  perPageInt,
			LastPage: lastPage,
			Data:     data,
		}

		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		encoded, err := json.Marshal(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = w.Write(encoded)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func main() {

	sqlxconn := "host=localhost port=5432 user=pack_admin password=packpass dbname=packform-db sslmode=disable"
	db, err := sqlx.Connect("postgres", sqlxconn)
	if err != nil {
		log.Fatalln(err)
	}

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalln(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalln(err)
	}

	mongodb := client.Database("packform-db")

	http.HandleFunc("/", generateHandler(db, mongodb))

	port := "8888"

	log.Infof("listening on :%s...", port)
	log.Fatalln(http.ListenAndServe(":"+port, nil))
}
