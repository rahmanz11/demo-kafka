package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
)

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
		GroupID:   "match-order-group",
		Topic:     "match-order",
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

			var data MatchedOrder
			// unmarshal json string to type data
			json.Unmarshal(m.Value, &data)

			// database connection string
			connStr := "postgresql://micros:micro$@localhost:5432/match_order_db?sslmode=disable"

			db, connerr := sql.Open("postgres", connStr)
			if connerr != nil {
				fmt.Printf("error while opening db con %s\n", connerr)
			} else {
				newId := 0
				sqlStatement := `INSERT INTO match_order_info (order_id, order_type, amt, _from, _to, pmt_method, sell_order_id) 
								VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`
				inserr := db.QueryRow(sqlStatement, data.OrderId, data.OrderType, data.Amt, data.From, data.To, data.PmtMethod, data.SellOrderId).Scan(&newId)

				if inserr != nil {
					fmt.Printf("error while inserting data into match_order_info table %s\n", inserr)
				} else {
					fmt.Println("New match-order record ID is:", newId)
					// when insertion is successful, there will be a valid id
					if newId > 0 {
						// kafka commit
						if comerr := r.CommitMessages(ctx, m); err != nil {
							fmt.Printf("failed to commit match-order messages: %s\n", comerr)
						} else {
							fmt.Println("committed match-order message")
						}
					}
				}
			}
		}
	}
}
