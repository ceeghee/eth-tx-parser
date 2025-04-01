package main

import (
	"eth-tx-parser/handler"
	"eth-tx-parser/parser"
	"log"
	"net/http"
)

func main() {
	parserInstance := parser.TransactionParser()
	go parserInstance.MonitorForLatestBlock()

	http.HandleFunc("/subscribe", handler.SubscribeHandler(parserInstance))
	http.HandleFunc("/transactions", handler.TransactionsHandler(parserInstance))
	http.HandleFunc("/current_block", handler.CurrentBlockHandler(parserInstance))

	log.Println("Server is live on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
