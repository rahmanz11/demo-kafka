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

type TransactionInit struct {
	TransactionId string    `json:"transaction_id"`
	CompletedAt   time.Time `json:"completed_at"`
}

// main function will be executed when this file is run
func main() {
	// initialize kafka connection and reader
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"172.31.128.1:29092"},
		GroupID:   "transactions-group",
		Topic:     "transactions",
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

			var transaction_init TransactionInit
			// unmarshal json string
			json.Unmarshal(m.Value, &transaction_init)

			// validate?

			// database connection string
			connStr := "postgresql://postgres:$@172.31.128.1:26257/transactions_db?sslmode=disable"

			db, connerr := sql.Open("postgres", connStr)
			if connerr != nil {
				fmt.Printf("Error while opening transactions db con %s\n", connerr)
			} else {
				sql_statement := `UPDATE transactions SET status = $1 WHERE transaction_id = $2;`
				res, err := db.Exec(sql_statement, "completed", transaction_init.TransactionId)

				if err != nil {
					fmt.Printf("Error occurred while updating transactions record. Reason: %s", err.Error())
				} else {

				}
				count, err := res.RowsAffected()
				if err != nil {
					fmt.Printf("Error occurred while getting rows affected of transactions. Reason: %s", err.Error())
				} else {
					if count > 0 {
						fmt.Printf("Transaction updated, id: %s", transaction_init.TransactionId)
					} else {
						fmt.Printf("No transaction row updated")
					}
				}
			}
		}
	}
}
