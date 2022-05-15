package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
)

type MatchedOrder struct {
	OrderId     string    `json:"orderId"`
	Amt         int32     `json:"amt"`
	From        string    `json:"from"`
	To          string    `json:"to"`
	PayWith     string    `json:"payWith"`
	PutProceeds string    `json:"putProceeds"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
}

type FundReceipt struct {
	FundReceiptId   string `json:"fundReceiptId"`
	TransactionType string `json:"transactionType"`
	CurrencyCode    string `json:"currencyCode"`
	CurrenyName     string `json:"currenyName"`
}

type MockRequest struct {
	ReceiptId   string `json:"receiptId"`
	FundAmt     int32  `json:"fundAmt"`
	PayWith     string `json:"payWith"`
	PutProceeds string `json:"putProceeds"`
}

type MockResponse struct {
	ReceiptId   string      `json:"receiptId"`
	FundReceipt FundReceipt `json:"fundReceipt"`
	PaidAt      string      `json:"paidAt"`
	Status      string      `json:"status"`
}

type PaidOrder struct {
	ReceiptId   string      `json:"receiptId"`
	From        string      `json:"from"`
	To          string      `json:"to"`
	FundAmt     int32       `json:"fundAmt"`
	PayWith     string      `json:"payWith"`
	PutProceeds string      `json:"putProceeds"`
	CreatedAt   time.Time   `json:"createdAt"`
	FundReceipt FundReceipt `json:"fundReceipt"`
	PaidAt      string      `json:"paidAt"`
	Status      string      `json:"status"`
}

func main() {
	// Step 1. Consume topic "matched-order"
	// initialize kafka connection and reader
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:9092"},
		GroupID:   "pay-match-order-group",
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
			// Step 2. Copy the topic message into local variable
			// unmarshal json string to type data
			json.Unmarshal(m.Value, &data)

			// database connection string
			connStr := "postgresql://micros:micro$@localhost:5432/match_order_db?sslmode=disable"

			db, connerr := sql.Open("postgres", connStr)
			if connerr != nil {
				fmt.Printf("error while opening db con %s\n", connerr)
			} else {
				newId := 0
				sqlStatement := `INSERT INTO match_order_info (receipt_id, _from, _to, fund_amt, pay_with, put_proceeds, created_at, status) 
								VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;`
				inserr := db.QueryRow(sqlStatement, data.OrderId, data.From, data.To, data.Amt, data.PayWith, data.PutProceeds, data.CreatedAt, "MATCHED").Scan(&newId)

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

						// Step 3. Make a request to the mock server
						mock_response := send_mock_request(data.OrderId, data.Amt, data.PayWith, data.PutProceeds)

						// Step 4. Produce a new topc "paid order" that includes mock server response
						publish_paid_order(mock_response, data, newId, db)

						// Step 5. Update matched_order_db record
						update_matched_order_paid(newId, db, mock_response.FundReceipt)

					}
				}
			}
		}
	}
}

func send_mock_request(receipt_id string, fund_amt int32, pay_with string, put_proceeds string) MockResponse {
	mock_request := MockRequest{}
	mock_response := MockResponse{}
	mock_request.ReceiptId = receipt_id
	mock_request.FundAmt = fund_amt
	mock_request.PayWith = pay_with
	mock_request.PutProceeds = put_proceeds
	body, _ := json.Marshal(&mock_request)

	// save match order data in database by calling the POST API for match order
	response, err := http.Post("https://1ea94e89-d393-4da6-a032-0c8c48faa311.mock.pstmn.io/fund", "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("cannot send mock request %s\n", err.Error())
	} else {
		response_data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("error while sending mock request: %s\n", err)
		} else {
			fmt.Printf("mock response data: %s\n", response_data)
			err := json.Unmarshal(response_data, &mock_response)
			if err != nil {
				fmt.Printf("error in mock response data unmarshal %s\n", err.Error())
			}
		}
	}

	return mock_response
}

func publish_paid_order(mock_response MockResponse, matched_order MatchedOrder, match_order_info_id int, db *sql.DB) bool {

	published := false
	// get kafka tcp connection -> broker address, topic name and kafka partition
	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", "paid-order", 0)
	if err != nil {
		fmt.Printf("failed to dial leader: %s\n", err)
	}

	if conn != nil {
		paid_order := PaidOrder{}
		paid_order.ReceiptId = matched_order.OrderId
		paid_order.From = matched_order.From
		paid_order.To = matched_order.To
		paid_order.FundAmt = matched_order.Amt
		paid_order.PayWith = matched_order.PayWith
		paid_order.PutProceeds = matched_order.PutProceeds
		paid_order.CreatedAt = matched_order.CreatedAt
		paid_order.FundReceipt = mock_response.FundReceipt
		paid_order.PaidAt = mock_response.PaidAt
		paid_order.Status = mock_response.Status

		// convert newOrder object to json string before publish message
		msg, _ := json.Marshal(paid_order)
		if msg != nil {
			_, err = conn.WriteMessages(
				kafka.Message{Value: msg},
			)
			if err != nil {
				fmt.Printf("failed to write match-order messages: %s\n", err)
			} else {
				fmt.Println("Paid order")
				published = true
			}
		}
	}

	return published
}

func update_matched_order_paid(match_order_info_id int, db *sql.DB, fund_receipt FundReceipt) {
	if fund_receipt.FundReceiptId != "" {
		sql_statement := `UPDATE match_order_info SET status = $1, fund_receipt_id = $2, transaction_type = $3, currency_code = $4, curreyncy_name = $5, paid_at = $6 WHERE id = $7;`
		res, err := db.Exec(sql_statement, "PAID", fund_receipt.FundReceiptId, fund_receipt.TransactionType, fund_receipt.CurrencyCode, fund_receipt.CurrenyName, time.Now(), match_order_info_id)
		if err != nil {
			fmt.Printf("error while updating match order info record-1 %s\n", err)
		} else {
			fmt.Println("match order info record updated, id-1:")
			fmt.Println(match_order_info_id)
		}
		count, err := res.RowsAffected()
		if err != nil {
			fmt.Printf("error while getting information of match order info rows updated-1 %s\n", err)
		} else {
			fmt.Printf("no. of match order info rows updated are-1:")
			fmt.Println(count)
		}
	} else {
		sql_statement := `UPDATE match_order_info SET status = $1, paid_at = $2 WHERE id = $3;`
		res, err := db.Exec(sql_statement, "PAID", time.Now(), match_order_info_id)
		if err != nil {
			fmt.Printf("error while updating match order info record-2 %s\n", err)
		} else {
			fmt.Println("match order info record updated, id-2:")
			fmt.Println(match_order_info_id)
		}
		count, err := res.RowsAffected()
		if err != nil {
			fmt.Printf("error while getting information of match order info rows updated-2 %s\n", err)
		} else {
			fmt.Printf("no. of match order info rows updated are-2:")
			fmt.Println(count)
		}
	}
}