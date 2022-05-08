package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
)

type Order struct {
	OrderId   string `json:"orderId"`
	OrderType string `json:"orderType"`
	Amt       int32  `json:"amt"`
	From      string `json:"from"`
	To        string `json:"to"`
	PmtMethod string `json:"pmtMethod"`
}

type MatchedOrder struct {
	OrderId     string `json:"orderId"`
	OrderType   string `json:"orderType"`
	Amt         int32  `json:"amt"`
	From        string `json:"from"`
	To          string `json:"to"`
	PmtMethod   string `json:"pmtMethod"`
	SellOrderId string `json:"sellOrderId"`
}

func main() {
	// initialize kafka connection and reader
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:9092"},
		GroupID:   "order-group",
		Topic:     "order",
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})

	if r != nil {

		ctx := context.Background()

		// run infinitely and fetch messages when available
		for {
			m, err := r.FetchMessage(ctx)
			if err != nil {
				fmt.Printf("error in fetch msg %s\n", err)
				break
			}

			var data Order
			// unmarshal json string to type data
			json.Unmarshal(m.Value, &data)

			// database connection string
			connStr := "postgresql://micros:micro$@localhost:5432/order_db?sslmode=disable"

			db, connerr := sql.Open("postgres", connStr)
			if connerr != nil {
				fmt.Printf("error while opening db con %s\n", connerr)
			} else {
				newId := 0
				sqlStatement := `INSERT INTO order_info (order_id, order_type, amt, _from, _to, pmt_method) 
								VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;`
				inserr := db.QueryRow(sqlStatement, data.OrderId, data.OrderType, data.Amt, data.From, data.To, data.PmtMethod).Scan(&newId)

				if inserr != nil {
					fmt.Printf("error while inserting data into order_info table %s\n", inserr)
				} else {
					fmt.Println("New order record ID is:", newId)

					// when insertion is successful, there will be a valid id
					if newId > 0 {
						matched := false

						if data.OrderType != "" && data.OrderType == "buy" {
							// go match incoming "buy" data
							matched = match(data, db)
						}

						if matched {
							// match successful then do kafka commit
							if comerr := r.CommitMessages(ctx, m); err != nil {
								fmt.Printf("failed to commit order messages: %s\n", comerr)
							} else {
								fmt.Println("committed order message successfully")
							}
						}
					}
				}
			}
		}
	}
}

// the match function
func match(data Order, db *sql.DB) bool {
	var id int
	var orderId string
	var from string

	success := false

	sqlStatement := `SELECT id, order_id, _from FROM order_info WHERE order_type = $1 AND amt = $2 ORDER BY created_at DESC LIMIT 1;`
	row := db.QueryRow(sqlStatement, "sell", data.Amt)

	switch err := row.Scan(&id, &orderId, &from); err {
	case sql.ErrNoRows:
		fmt.Printf("No sell order amt matched for order-id: %s\n", data.OrderId)
	case nil:
		fmt.Printf("Sell order amt matched for order-id: %s\n", data.OrderId)

		sqlStatement = `UPDATE order_info SET matched = $1 WHERE id = $2;`
		res, err := db.Exec(sqlStatement, true, id)

		if err != nil {
			fmt.Printf("error while updating order record %s\n", err)
		} else {
			fmt.Println("order record updated, id:")
			fmt.Println(id)
		}
		count, err := res.RowsAffected()
		if err != nil {
			fmt.Printf("error while getting information of order rows updated %s\n", err)
		} else {
			fmt.Printf("no. of order rows updated are:")
			fmt.Println(count)
		}

		var matchOrder MatchedOrder

		matchOrder.OrderId = data.OrderId
		matchOrder.OrderType = data.OrderType
		matchOrder.Amt = data.Amt
		matchOrder.From = from
		matchOrder.To = data.To
		matchOrder.PmtMethod = data.PmtMethod
		matchOrder.SellOrderId = orderId

		body, _ := json.Marshal(&matchOrder)

		// save match order data in database by calling the POST API for match order
		response, err := http.Post("http://188.166.32.136:90/api/v2/match-order", "application/json", bytes.NewBuffer(body))
		if err != nil {
			fmt.Printf("cannot send match-order api request %s\n", err.Error())
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("error while posting match-order to api: %s\n", err)
			success = true
		} else {
			json.MarshalIndent(responseData, "", "\t")
		}
	default:
		fmt.Printf("error while opening db con %s\n", err)
	}

	return success
}
