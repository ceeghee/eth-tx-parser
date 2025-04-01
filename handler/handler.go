package handler

import (
	"encoding/json"
	"eth-tx-parser/parser"
	"fmt"
	"net/http"
)

func SubscribeHandler(tp *parser.TxParser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		addr := r.URL.Query().Get("address")
		if addr == "" {
			http.Error(w, "Address is required", http.StatusBadRequest)
			return
		}
		if tp.Subscribe(addr) {
			fmt.Fprintf(w, "Address Successfully subscribed: %s", addr)
		} else {
			fmt.Fprintf(w, "Address is already subscribed!")
		}
	}
}


func TransactionsHandler(tp *parser.TxParser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		addr := r.URL.Query().Get("address")
		if addr == "" {
			http.Error(w, "Address parameter is missing", http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(tp.Transactions(addr))
	}
}

func CurrentBlockHandler(tp *parser.TxParser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%d", tp.CurrentBlock())
	}
}