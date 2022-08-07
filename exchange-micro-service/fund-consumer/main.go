package main

import (

	// for printing logs

	// for http usage

	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	// for gin framework
	kafka "github.com/segmentio/kafka-go"

	// for postgres library
	_ "github.com/lib/pq"
	// for kafka go library
)

type ExchangeOrder struct {
	Blank         string    `json:"blank"`
	From          string    `json:"from"`
	Fund          string    `json:"fund"`
	Amt           float64   `json:"amt"`
	Re            string    `json:"re"`
	TransactionId string    `json:"transaction_id"`
	ReceivedAt    time.Time `json:"received_at"`
	ValidatedAt   time.Time `json:"validated_at"`
	SubmittedAt   time.Time `json:"submitted_at"`
}

type TransactionInit struct {
	TransactionId string    `json:"transaction_id"`
	CompletedAt   time.Time `json:"completed_at"`
}

// main function will be executed when this file is run
func main() {
	// initialize kafka connection and reader
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"172.31.128.1:29092"},
		GroupID:   "fund-group",
		Topic:     "exchange_orders",
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

			var exchange_order ExchangeOrder
			// unmarshal json string
			json.Unmarshal(m.Value, &exchange_order)

			// validate?

			// database connection string
			connStr := "postgresql://postgres:$@172.31.128.1:26257/transactions_db?sslmode=disable"

			db, connerr := sql.Open("postgres", connStr)
			if connerr != nil {
				fmt.Printf("Error while opening transactions database connection. Reason: %s", connerr)
			} else {
				sqlStatement := `INSERT INTO transactions (transaction_id, _from, fund, amt, re, received_at, 
								validated_at, submitted_at, completed_at, status, system_notes) 
								VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
				completed_at := time.Now()
				_, err := db.Exec(sqlStatement, exchange_order.TransactionId, exchange_order.From, exchange_order.Amt,
					exchange_order.Amt, exchange_order.Re, exchange_order.ReceivedAt, exchange_order.ValidatedAt,
					exchange_order.SubmittedAt, completed_at, "pending", "default note")

				if err != nil {
					fmt.Printf("Error while inserting data into transactions table. Reason: %s", err.Error())
				} else {
					// produce to transactions
					var transaction_init TransactionInit
					transaction_init.TransactionId = exchange_order.TransactionId
					transaction_init.CompletedAt = completed_at

					conn, err := kafka.DialLeader(context.Background(), "tcp", "172.31.128.1:29092", "transactions", 0)
					if err != nil {
						fmt.Printf("Failed to dial leader: %s", err)
					}
					// convert ExchangeOrder object to json string before publish message
					msg, _ := json.Marshal(transaction_init)
					if msg != nil {
						_, err = conn.WriteMessages(
							kafka.Message{Value: msg},
						)
						if err != nil {
							fmt.Printf("Failed to produce transaction message. Reason: %s", err)
						}
					}
				}
			}
		}
	}
}
