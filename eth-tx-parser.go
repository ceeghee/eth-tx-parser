package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"sync"
	"time"
)

const ethRPCURL = "https://ethereum-rpc.publicnode.com"

type Transaction struct {
	Sender string `json:"from"`
	Amount string `json:"value"`
	Receiver string `json:"to"`
	TxHash string `json:"hash"`
}

type RPCRequest struct {
	Method  string        `json:"method"`
	JSONRPC string        `json:"jsonrpc"`
	ID      int           `json:"id"`
	Params  []interface{} `json:"params"`
}

type RPCResponse struct {
	ID      int    `json:"id"`
	Result  string `json:"result"`
	JSONRPC string `json:"jsonrpc"`
}

type TxParser struct {
	lock         sync.Mutex
	latestBlock  int
	subscribers  map[string]bool
	txRecords    map[string][]Transaction
}

func (tp *TxParser) monitorFoLatestBlock() {
	for {
		newBlock := tp.fetchLatestBlock()
		if newBlock > tp.latestBlock {
			log.Printf("New block detected: %d\n", newBlock)
			tp.lock.Lock()
			tp.latestBlock = newBlock
			tp.lock.Unlock()
		}
		time.Sleep(10 * time.Second)
	}
}

func TransactionParser() *TxParser {
	return &TxParser{
		latestBlock: 0,
		subscribers: make(map[string]bool),
		txRecords:   make(map[string][]Transaction),
	}
}

func (tp *TxParser) CurrentBlock() int {
	tp.lock.Lock()
	defer tp.lock.Unlock()
	return tp.latestBlock
}


func (tp *TxParser) Transactions(addr string) []Transaction {
	tp.lock.Lock()
	defer tp.lock.Unlock()
	return tp.txRecords[addr]
}


func (tp *TxParser) Subscribe(addr string) bool {
	tp.lock.Lock()
	defer tp.lock.Unlock()
	if tp.subscribers[addr] {
		return false
	}
	tp.subscribers[addr] = true
	return true
}

func (tp *TxParser) fetchLatestBlock() int {
	request := RPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "eth_blockNumber",
		Params:  []interface{}{},
	}
	response, err := sendRPCRequest(request)
	if err != nil {
		log.Println("Unabled to fetch block:", err)
		return tp.latestBlock
	}
	blockNum, _ := new(big.Int).SetString(response.Result[2:], 16)
	return int(blockNum.Int64())
}

func sendRPCRequest(req RPCRequest) (*RPCResponse, error) {
	payload, _ := json.Marshal(req)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(ethRPCURL, "application/json", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var rpcResp RPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
		return nil, err
	}
	return &rpcResp, nil
}

func main() {
	parser := TransactionParser()
	go parser.monitorFoLatestBlock()

	http.HandleFunc("/subscribe", func(w http.ResponseWriter, r *http.Request) {
		addr := r.URL.Query().Get("address")
		if addr == "" {
			http.Error(w, "Address is required", http.StatusBadRequest)
			return
		}
		if parser.Subscribe(addr) {
			fmt.Fprintf(w, "Address Successfully subscribed: %s", addr)
		} else {
			fmt.Fprintf(w, "Address is already subscribed!")
		}
	})

	http.HandleFunc("/current_block", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%d", parser.CurrentBlock())
	})

	http.HandleFunc("/transactions", func(w http.ResponseWriter, r *http.Request) {
		addr := r.URL.Query().Get("address")
		if addr == "" {
			http.Error(w, "Address is missing", http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(parser.Transactions(addr))
	})

	log.Println("Server running on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
