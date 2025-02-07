# Ethereum Blockchain Parser

A Go-based Ethereum blockchain parser that allows real-time monitoring of transactions for subscribed addresses.

## Features

- Real-time Ethereum blockchain parsing
- Address subscription management
- Transaction monitoring for subscribed addresses
- REST API for interaction
- In-memory storage (easily extendable)
- Thread-safe operations
- Comprehensive logging

## Pre-requisites

- Go 1.19 or higher
- Access to an Ethereum node (default: https://ethereum-rpc.publicnode.com)

## Project Structure

The project follows a clean architecture pattern with the following structure:
```
   eth-parser/
    ├── README.md
    ├── go.mod
    ├── Makefile
    ├── cmd/
    │   └── main.go                    # Application entry point
    ├── internal/
    │   ├── api/
    │   │   ├── server.go             # HTTP API implementation
    │   │   └── server_test.go        # API tests
    │   ├── parser/
    │   │   ├── parser.go             # Core parser implementation
    │   │   └── parser_test.go        # Parser unit tests
    │   ├── rpc/
    │   │   ├── client.go             # Ethereum JSON-RPC client
    │   │   └── client_test.go        # RPC client tests
    │   └── storage/
    │       ├── memory.go             # In-memory storage implementation
    │       └── memory_test.go        # Storage tests
    └── pkg/
        └── types/
            └── types.go              # Shared types and interfaces
```

## Installation

1. Clone the repository:
```bash
git clone https://github.com/tanryberdi/ethparser
cd ethparser
```

2. Install the dependencies:
```bash
make deps
```

3. Build the application:
```bash
make build
```

4. Run the application:
```bash
make run
```

By default, the service runs on port 8080. You can modify the port by setting the PORT environment variable:
```bash
PORT=3000 ./ethparser
```

5. Testing

Run the tests using the following command:

```bash
make test
```

## API Endpoints

1. Subscribe to an address:

Subscribe to receive notifications for a specific Ethereum address.

```bash
curl -X POST http://localhost:8080/subscribe \
  -H "Content-Type: application/json" \
  -d '{
    "address": "0x742d35Cc6634C0532925a3b844Bc454e4438f44e"
  }'
```

Response:

```json
{
  "success": true,
  "message": "Address successfully subscribed"
}
```

2. Get Transactions

Retrieve all transactions for a subscribed address.

```bash
curl -X GET "http://localhost:8080/transactions?address=0x742d35Cc6634C0532925a3b844Bc454e4438f44e"
```

Response:

```json
{
  "address": "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
  "transactions": [
    {
      "hash": "0x...",
      "from": "0x...",
      "to": "0x...",
      "value": "0x...",
      "blockNumber": 14000000,
      "timestamp": 1632150000
    }
  ]
}
```

3. Get current block

Get the last parsed block number.

```bash
curl -X GET http://localhost:8080/current-block
```

Response:

```json
{
  "blockNumber": 14000000
}
```

