package main

import (
	"fmt"
	"strings"
)

func buildQuery(filters []string, idx int) (query string, where string) {

	where = ""
	if len(filters) > 0 {
		where = "where " + strings.Join(filters, " and ")
	}

	query = `with some_orders as (
				select id, order_name, created_at, customer_id
				from orders
				` + where + ` order by created_at desc
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
			order by created_at desc`
	return
}
