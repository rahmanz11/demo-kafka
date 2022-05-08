package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
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
	r := gin.Default()
	r.POST("/api/v2/order", PostOrder)
	r.GET("/api/v2/list", GetOrder)
	r.POST("/api/v2/match-order", PostMatchOrder)
	r.Run("localhost:89")
}

func PostOrder(c *gin.Context) {
	var newOrder Order

	// call BindJSON to bind the received JSON to newOrder
	if err := c.BindJSON(&newOrder); err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	published := false

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", "order", 0)
	if err != nil {
		fmt.Printf("failed to dial leader: %s\n", err)
	}

	if conn != nil {
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

func PostMatchOrder(c *gin.Context) {
	var newMatchOrder MatchedOrder

	// call BindJSON to bind the received JSON to newMatchOrder
	if err := c.BindJSON(&newMatchOrder); err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	published := false

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", "match-order", 0)
	if err != nil {
		fmt.Printf("failed to dial leader: %s\n", err)
	}

	if conn != nil {
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

func GetOrder(c *gin.Context) {

	orders := make([]*Order, 0)

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
		defer rows.Close()
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
				//TODO insert into list
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
}
