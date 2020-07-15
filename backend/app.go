package main

import (
	"log"
	"net/http"

	"encoding/json"
)

type HTTPResponse struct {
	CurrentPage int     `json:"current_page,omitempty"`
	Total       int     `json:"total,omitempty"`
	From        int     `json:"from,omitempty"`
	To          int     `json:"to,omitempty"`
	NextPageURL string  `json:"next_page_url,omitempty"`
	PrevPageURL string  `json:"prev_page_url,omitempty"`
	PerPage     int     `json:"per_page,omitempty"`
	LastPage    int     `json:"last_page,omitempty"`
	Data        []Order `json:"data,omitempty"`
}

type Order struct {
	OrderName       string `json:"order_name,omitempty"`
	CustomerCompany string `json:"customer_company,omitempty"`
	CustomerName    string `json:"customer_name,omitempty"`
	OrderDate       string `json:"order_date,omitempty"`
	DeliveredAmount string `json:"delivered_amount,omitempty"`
	TotalAmount     string `json:"total_amount,omitempty"`
}

func handler(w http.ResponseWriter, r *http.Request) {

	// create dummy response
	resp := HTTPResponse{
		CurrentPage: 1,
		Total:       1,
		From:        1,
		To:          1,
		NextPageURL: "",
		PrevPageURL: "",
		PerPage:     1,
		LastPage:    1,
		Data: []Order{
			{
				OrderName:       "C19190 Christmas",
				CustomerCompany: "Sony Ericsson",
				CustomerName:    "Dr. Harold Senger",
				OrderDate:       "2020-01-02T15:34:12Z",
				DeliveredAmount: "$99.11",
				TotalAmount:     "$99.11",
			},
		},
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
}

func main() {
	http.HandleFunc("/", handler)

	log.Println("listening...")
	log.Fatal(http.ListenAndServe(":8888", nil))
}
