package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	// for printing logs
	"fmt"

	// for http usage
	"net/http"

	// for gin framework
	"github.com/gin-gonic/gin"

	// for postgres library
	_ "github.com/lib/pq"

	// for kafka go library
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
	OrderId     string    `json:"orderId"`
	Amt         int32     `json:"amt"`
	From        string    `json:"from"`
	To          string    `json:"to"`
	PayWith     string    `json:"payWith"`
	PutProceeds string    `json:"putProceeds"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	SellOrderId string    `json:"sellOrderId"`
	OrderType   string    `json:"orderType"`
}

// main function will be executed when this file is run
func main() {

	// gin framework for REST API
	r := gin.Default()

	// API endpoints
	r.POST("/api/v2/order", PostOrder)
	r.GET("/api/v2/list", GetOrder)
	r.POST("/api/v2/match-order", PostMatchOrder)

	// API will run at mentioned address
	r.Run("188.166.32.136:90")
}

// handler function for Order Post
func PostOrder(c *gin.Context) {
	var newOrder Order

	// call BindJSON to bind the received JSON to newOrder
	if err := c.BindJSON(&newOrder); err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	published := false

	// get kafka tcp connection -> broker address, topic name and kafka partition
	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", "order", 0)
	if err != nil {
		fmt.Printf("failed to dial leader: %s\n", err)
	}

	if conn != nil {
		// convert newOrder object to json string before publish message
		msg, _ := json.Marshal(newOrder)
		if msg != nil {
			_, err = conn.WriteMessages(
				kafka.Message{Value: msg},
			)
			if err != nil {
				fmt.Printf("failed to write order messages: %s\n", err)
			} else {
				published = true
			}
		}
	}

	if published {
		// return success
		c.IndentedJSON(http.StatusCreated, http.StatusCreated)
	} else {
		// return bad gateway
		c.IndentedJSON(http.StatusBadGateway, http.StatusBadGateway)
	}

}

// handler function for MatchOrder Post
func PostMatchOrder(c *gin.Context) {
	var newMatchOrder MatchedOrder

	// call BindJSON to bind the received JSON to newMatchOrder
	if err := c.BindJSON(&newMatchOrder); err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	published := false

	// get kafka tcp connection -> broker address, topic name and kafka partition
	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", "match-order", 0)
	if err != nil {
		fmt.Printf("failed to dial leader: %s\n", err)
	}

	if conn != nil {
		// convert newOrder object to json string before publish message
		msg, _ := json.Marshal(newMatchOrder)
		if msg != nil {
			_, err = conn.WriteMessages(
				kafka.Message{Value: msg},
			)
			if err != nil {
				fmt.Printf("failed to write match-order messages: %s\n", err)
			} else {
				published = true
			}
		}
	}

	if published {
		// return success
		c.IndentedJSON(http.StatusCreated, http.StatusCreated)
	} else {
		// return bad gateway
		c.IndentedJSON(http.StatusBadGateway, http.StatusBadGateway)
	}

}

// handler function for Order List
func GetOrder(c *gin.Context) {

	// order list
	orders := make([]*Order, 0)

	// database connection string
	connStr := "postgresql://micros:micro$@localhost:5432/order_db?sslmode=disable"

	db, connerr := sql.Open("postgres", connStr)
	if connerr != nil {
		fmt.Printf("error while opening db con %s\n", connerr)
	} else {

		sqlStatement := `SELECT order_id, order_type, amt, _from, _to, pmt_method FROM order_info WHERE order_type = $1 AND matched = $2;`
		rows, err := db.Query(sqlStatement, "sell", false)

		if err != nil {
			// handle this error better than this
			panic(err)
		}

		// close the resultset when iteration is done
		defer rows.Close()

		// iterate each row
		for rows.Next() {
			var orderId string
			var orderType string
			var amt int32
			var from string
			var to string
			var pmtMethod string

			switch err := rows.Scan(&orderId, &orderType, &amt, &from, &to, &pmtMethod); err {
			case sql.ErrNoRows:
				fmt.Println("No sell order with matched false available")
			case nil:
				fmt.Println("Sell order available with matched is false")
				var order Order
				order.OrderId = orderId
				order.OrderType = orderType
				order.Amt = amt
				order.From = from
				order.To = to
				order.PmtMethod = pmtMethod

				orders = append(orders, &order)

			default:
				fmt.Printf("error while opening db conn %s\n", err)
			}
		}

		// get any error encountered during iteration
		err = rows.Err()
		if err != nil {
			fmt.Printf("error while fetching rows %s\n", err)
		}
	}

	// response writer
	c.IndentedJSON(http.StatusOK, orders)
}
