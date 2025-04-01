# Ethereum Transaction Parser

## Overview
This project implements an Ethereum transaction parser that monitors the blockchain for transactions involving subscribed addresses. It provides an HTTP API for subscribing to addresses, retrieving transaction history, and querying the latest processed block.

## Features
- Monitors Ethereum blockchain for new blocks.
- Allows users to subscribe to Ethereum addresses.
- Stores transaction history in memory (extensible to external databases).
- Exposes an HTTP API for interaction.

## Setup Instructions
### Prerequisites
- Go 1.18+
- Internet connection for Ethereum RPC access

### Installation
1. Clone the repository:
   ```sh
   git clone https://github.com/ceeghee/eth-tx-parser.git
   cd eth-tx-parser
   ```
2. Run the application:
   ```sh
   go run eth-tx-parser.go
   ```
3. The server starts on `http://localhost:8080`

## API Endpoints
### Subscribe to an Address
**Request:**
```sh
curl -X GET "http://localhost:8080/subscribe?address=sampleEthAddress"
```
**Response:**
```sh
Address Successfully subscribed sampleEthAddress
```

### Get Transactions for an Address
**Request:**
```sh
curl -X GET "http://localhost:8080/transactions?address=0xYourEthereumAddress"
```
**Response:**
```json
[
  {
		"from": "senderEthAddress",
    "value": "0.0004",
    "hash": "transaction hash",
    "to": "receiverEthAddress",
  }
]
```

### Get Latest Processed Block
**Request:**
```sh
curl -X GET "http://localhost:8080/current_block"
```
**Response:**
```sh
12345678
```

## Implementation Details
### 1. **Ethereum JSON-RPC Communication**
- Uses `eth_blockNumber` to fetch the latest block.
- Utilizes `sendRPCRequest` to abstract RPC calls.
- Avoids external dependencies by using Go’s `net/http` and `encoding/json`.

### 2. **Memory Storage with Sync Mechanism**
- Transactions and subscriptions are stored in memory for quick access.
- `sync.Mutex` ensures safe concurrent access.
- Can be extended to support persistent storage.

### 3. **Efficient Blockchain Monitoring**
- A goroutine continuously fetches new blocks.
- Only updates the latest block if a new one is detected.
- Uses `time.Sleep(10 * time.Second)` to avoid unnecessary requests.

## Future Enhancements
- Implement WebSockets for real-time transaction monitoring.
- Integrate a database for persistent storage preferably a nosql database
- Add support for filtering transactions based on type (incoming/outgoing).



