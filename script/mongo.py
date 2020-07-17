import csv
from datetime import datetime
from pymongo import MongoClient

# This assumes database and tables are created.
if __name__ == "__main__":
    client = MongoClient('mongodb://localhost:27017')
    db = client['packform-db']

    tasks = {}
    tasks["1"] = {
        "csv_name": "../data/mongo_customer_companies.csv",
        "collname": "customer_companies",
        "objfunc": lambda row: {
            "company_id": int(row[0]),
            "company_name": row[1]
        }
    }
    tasks["2"] = {
        "csv_name": "../data/mongo_customers.csv",
        "collname": "customers",
        "objfunc": lambda row: {
            "user_id": row[0],
            "login": row[1],
            "password": row[2],
            "name": row[3],
            "company_id": int(row[4]),
            "credit_cards": row[5]
        }
    }

    for _, v in tasks.items():
        print(v)
        with open(v['csv_name']) as csv_file:
            csv_reader = csv.reader(csv_file, delimiter=',')
            line_count = 0
            for row in csv_reader:
                if line_count == 0:
                    line_count += 1 # skip header
                    continue
                obj = v['objfunc'](row)
                try:
                    _ = db[v['collname']].insert_one(obj)
                except Exception as e:
                    print("unable to execute insert for obj: {}".format(obj))
                    break
                print("success insert obj: {}".format(obj))
    print("ok")