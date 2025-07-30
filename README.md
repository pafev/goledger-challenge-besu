# GoLedger Besu Challenge

## Project Description

Personal implementation (Paulo Fernando Vilarim) for the GoLedger technical challenge, which consists of creating a Go application that interacts with a Hyperledger Besu QBFT blockchain network. The application must be able to interact with smart contracts deployed on the Besu network, manage smart contract variable values, synchronize data between the blockchain and an SQL database (Postgres in the current implementation), and expose functionalities through a REST API.

## Prerequisites

Ensure the following are installed:

* Go (version 1.24)
* Node.js and NPM (recommended via NVM)
* Docker and Docker Compose
* Hardhat
* Hyperledger Besu

## How to Run

### 1. Environment Setup

Generate the environment variables file from the template:

```bash
cp .env.example .env
```

### 2. Start the Besu Client

Run the following commands to start the local Besu network:

```bash
cd scripts/besu
chmod +x ./startDev.sh
./startDev.sh
cd ../..
```

This command will:

* Start a local Besu network with 4 nodes
* Deploy the SimpleStorage smart contract
* Display the contract address (note this address to configure in `.env`)

### 3. Start the Database

Run the script to start the PostgreSQL container:

```bash
chmod +x ./scripts/db/startDev.sh
./scripts/db/startDev.sh
```

This command will:

* Start a container running a PostgreSQL 1.17 database
* Display the database URL (automatically assigned to `.env`, but confirm it)
  Note: if you choose to run the database locally, ensure the connection URL in `.env` contains the query param `?sslmode=disable`. Verify this.

### 4. Configure Environment Variables

Fill in the `.env` file with the information gathered from the previous steps. It should look roughly like this:

```env
# Application info
APP_NAME=goledger-challenge-besu
APP_DOMAIN="localhost"
APP_ENV="development"
APP_PORT=5000

# Database configuration
DATABASE_URL="<your_database_connection_url>"

# Besu network settings
BESU_URL=http://localhost:8545
SMART_CONTRACT_ADDR="<deployed_contract_address>"
SMART_CONTRACT_ABI_PATH="scripts/besu/artifacts/contracts/SimpleStorage.sol/SimpleStorage.json"
```

### 5. Install Dependencies

Download Go dependencies:

```bash
go mod download
go mod tidy
```

### 6. Run the Application

Run the app:

```bash
go run cmd/main.go
```

The application will be available at `http://localhost:5000`.

## Features and Endpoints

### GET /api/v1/smart-contract/

* Retrieves the current value stored in the smart contract
* Returns JSON with the current value

### GET /api/v1/smart-contract/check-value/\:value

* Compares the smart contract value with the provided value
* Returns JSON indicating equality (`true`/`false`)

### POST /api/v1/smart-contract/set-value

* Sets a new value in the smart contract
* Request body (JSON):

```json
{
  "value": 42,
  "privateKey": "..."
}
```

* Returns JSON confirming the transaction

### POST /api/v1/smart-contract/sync

* Synchronizes the smart contract value with the PostgreSQL database
* Returns JSON confirming synchronization

## About the Technologies

### Hyperledger Besu

Hyperledger Besu is an open-source Ethereum client developed by ConsenSys and later donated to the Hyperledger Foundation. Besu implements the Ethereum Virtual Machine (EVM) and supports major Ethereum consensus algorithms, including Proof of Work (PoW) and Proof of Stake (PoS).

### Ethereum and Smart Contracts

Ethereum is a decentralized blockchain platform that enables the execution of smart contracts—self-executing programs with the terms of the agreement directly written in code, offering security and trust without the need for censorship, fraud, or third-party interference.

### Besu as an Ethereum Client

As an Ethereum client, Besu provides:

* Synchronization with the Ethereum network
* Transaction and block validation
* Smart contract execution
* Block mining/validation capabilities

For development and testing environments, Besu can be configured to run local private networks, allowing application development and testing without mainnet gas costs.

## Application Architecture

The application follows Clean Architecture principles, but avoids over-engineering due to the reduced project scope. It maintains modularity, applied design patterns, and proper error handling for scalability and maintainability. The project has a clear division between application and domain layers. The structure follows a feature-based separation within each layer.

### Application Layers

**Domain Layer (`domain/`)**

* Contains the business rules and core entities
* Implements the Repository pattern to decouple business logic from client implementations (Postgres DB and Ethereum Besu client)
* Organized by feature (e.g., `smart_contract/`)

**Application Layer (`app/`)**

* Coordinates domain logic with external interfaces (HTTP requests)
* Contains use cases and services
* Organized by feature, matching the domain layer folder structure

### Repository Pattern

The domain layer uses the Repository pattern to abstract interactions with external clients:

* **PostgreSQL**: For data persistence
* **Besu (Ethereum Client)**: For blockchain interaction

This pattern makes the business logic independent of external system implementation details, making testing and maintenance easier.

```
├── internal/domain/
│   └── smart_contract/
│       ├── repository-besu.go
│       ├── repository-db.go
│       ├── model.go
├── internal/app/
│   └── smart_contract/
│       ├── handler.go
│       └── service.go
└── scripts/
    ├── db/
    └── besu/
```

## Technologies Used

* **Language**: Go
* **Blockchain**: Hyperledger Besu
* **Ethereum Client Lib**: go-ethereum
* **Database**: PostgreSQL
* **API**: REST
* **Containerization**: Docker & Docker Compose
* **Smart Contract**: SimpleStorage

## Project Structure

```
.
├── cmd/
│   └── main.go
├── configs/
│   └── app/
│       ├── config.go
│   └── besu/
│       ├── config.go
│   └── db/
│       ├── config.go
│   └── http/
│       └── config.go
├── internal/
│   └── app/
│       └── smart_contract/
│           ├── handler.go
│           └── service.go
│   └── domain/
│       └── smart_contract/
│           ├── repository-besu.go
│           ├── repository-db.go
│           └── model.go
├── scripts/
│   ├── besu/
│   └── db/
├── .env.example
├── go.mod
├── go.sum
└── README.md
```

## Usage Examples

### Retrieve contract value

```bash
curl -X GET http://localhost:8080/api/v1/smart-contracts/
```

### Set new value

```bash
curl -X POST http://localhost:8080/api/v1/smart-contracts/set-value \
  -H "Content-Type: application/json" \
  -d '{
    "value": 123,
    "privateKey": "8f2a55949038a9610f50fb23b5883af3b4ecb3c3bb792cbcefbd1542c692be63"
  }'
```

### Check value

```bash
curl -X GET http://localhost:8080/api/v1/smart-contracts/check-value/123
```

### Sync with database

```bash
curl -X POST http://localhost:8080/api/v1/smart-contracts/sync
```

## Error Handling

The application implements comprehensive error handling for:

* Besu network connection failures
* Blockchain transaction errors
* Database connectivity issues
* Input validation errors
* Network operation timeouts
* Route parameter parsing
* Private key authentication

All errors return proper HTTP responses with descriptive messages or error logs.

## Technical Notes

### Smart Contract Interaction

* Write transactions: `bind.NewKeyedTransactorWithChainID`
* Read calls: `bind.CallOpts`
* ABI: Auto-loaded from Hardhat artifacts

### Security

* Private keys provided via requests (not stored)
* Sensitive data protected via environment variables
* ABI read from source files
* Input validation on all endpoints
* Blockchain operation timeouts
* Database connection handling

### Performance

* Database connection pooling
* Ethereum client reuse
* Optimized data structures
* Configured timeouts

### Logging

Structured logging for debugging and production monitoring.

## Contribution

This project is part of a GoLedger technical challenge, developed by me, Paulo Fernando Vilarim (pafev).

**Developed for the GoLedger Challenge**
