import psycopg2
import csv
from datetime import datetime

# This assumes database and tables are created.
if __name__ == "__main__":
    try:
        conn = psycopg2.connect("dbname='packform-db' user='pack_admin' host='localhost' password='packpass'")
    except Exception as e:
        print "unable to connect to database {}".format(e)

    cur = conn.cursor()

    tasks = {}
    tasks["delivery"] = {
        "csv_name": "../data/pg_deliveries.csv",
        "cmdfunc": lambda row: "insert into delivery (id,order_item_id,delivered_quantity) values ({},{},{})".format(int(row[0]), int(row[1]), int(row[2]))
    }
    tasks["order_items"] = {
        "csv_name": "../data/pg_order_items.csv",
        "cmdfunc": lambda row: "insert into order_items (id,order_id,price_per_unit,quantity,product) values ({},{},{},{},'{}')".format(int(row[0]), int(row[1]), float(row[2]), int(row[3]), row[4])
    }
    tasks["orders"] = {
        "csv_name": "../data/pg_orders.csv",
        "cmdfunc": lambda row: "insert into orders (id,created_at,order_name,customer_id) values ({},'{}','{}','{}')".format(int(row[0]), datetime.strptime(row[1], '%Y-%m-%dT%H:%M:%SZ'), row[2], row[3])
    }

    for k, v in tasks.items():
        # print(k, v)
        with open(v['csv_name']) as csv_file:
            csv_reader = csv.reader(csv_file, delimiter=',')
            line_count = 0
            for row in csv_reader:
                if line_count == 0:
                    line_count += 1 # skip header
                    continue
                print(row)
                cmd = v['cmdfunc'](row)
                try:
                    cur.execute(cmd)
                except Exception as e:
                    print("unable to execute {}: {}".format(cmd, e))
                print("success exec: {}".format(cmd))

    conn.commit()

    cur.close()
    conn.close()
    print "ok"