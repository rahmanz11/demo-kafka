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
	OrderType   string    `json:"orderType"`
	Amt         int32     `json:"amt"`
	From        string    `json:"from"`
	To          string    `json:"to"`
	PmtMethod   string    `json:"pmtMethod"`
	SellOrderId string    `json:"sellOrderId"`
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
				sqlStatement := `INSERT INTO match_order_info (order_id, order_type, amt, _from, _to, pmt_method, sell_order_id, pay_with, put_proceeds, status, created_at) 
								VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id;`
				inserr := db.QueryRow(sqlStatement, data.OrderId, data.OrderType, data.Amt, data.From, data.To, data.PmtMethod, data.SellOrderId, data.PayWith, data.PutProceeds, data.Status, data.CreatedAt).Scan(&newId)

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

						mock_response := send_mock_request(data.OrderId, data.Amt, data.PayWith, data.PutProceeds)

						if mock_response.Status != "" && mock_response.Status == "PAID" {
							publish_paid_order(mock_response, data, newId, db)
						}
					}
				}
			}
		}
	}
}

func publish_paid_order(mock_response MockResponse, matched_order MatchedOrder, match_order_info_id int, db *sql.DB) {
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
				mark_matched_order_paid(match_order_info_id, db)
			}
		}
	}
}

func mark_matched_order_paid(match_order_info_id int, db *sql.DB) {
	sql_statement := `UPDATE match_order_info SET paid = $1 WHERE id = $2;`
	res, err := db.Exec(sql_statement, true, match_order_info_id)

	if err != nil {
		fmt.Printf("error while updating match order info record %s\n", err)
	} else {
		fmt.Println("match order info record updated, id:")
		fmt.Println(match_order_info_id)
	}
	count, err := res.RowsAffected()
	if err != nil {
		fmt.Printf("error while getting information of match order info rows updated %s\n", err)
	} else {
		fmt.Printf("no. of match order info rows updated are:")
		fmt.Println(count)
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
		fmt.Printf("cannot send match-order api request %s\n", err.Error())
	}

	response_data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("error while posting match-order to api: %s\n", err)
	} else {
		err := json.Unmarshal(response_data, &mock_response)
		if err != nil {
			fmt.Printf("error in mock request %s\n", err.Error())
		}
	}

	return mock_response
}
