# EventFlow

EventFlow is a real-time data streaming platform designed for efficient data ingestion and storage using RabbitMQ, PostgreSQL, and Redis.

## 🚀 Features
- **Producer Service**: Reads user data from CSV and sends it to RabbitMQ.
- **Consumer Service**: Processes messages from RabbitMQ and writes to PostgreSQL & Redis.

## 🛠 Tech Stack
- **Backend**: Golang, RabbitMQ, PostgreSQL, Redis
- **Infrastructure**: Docker, Docker Compose

## 🔧 Setup & Installation

### 1️⃣ Clone Repository
```sh
git clone https://github.com/vignesh-pandey/event-flow.git
cd event-flow
```

## API Documentation
Find API docs on [Postman](https://documenter.getpostman.com/view/28322535/2sAYdfqWRD).

# Producer Service

Handles CSV data extraction and publishing to RabbitMQ.

## Prerequisites
- **Go** ([Install](https://go.dev/doc/install))
- **Docker & Docker Compose** ([Install](https://docs.docker.com/get-docker/))
- **RabbitMQ** ([Download](https://www.rabbitmq.com/download.html))
- **cURL** (for API testing)

### Install Go
#### Windows
1. Download from [Go Downloads](https://go.dev/dl/).
2. Verify:
   ```sh
   go version
   ```
#### macOS (Homebrew)
```sh
brew install go
```
#### Linux (Ubuntu/Debian)
```sh
sudo apt update && sudo apt install -y golang
```

## Run Producer Locally

1. Navigate to `producer-service`:
   ```sh
   cd producer-service
   ```
2. Install dependencies:
   ```sh
   go mod tidy
   ```
3. Ensure RabbitMQ is running.
4. Create `config.json`:
   ```json
   {
     "rabbitmq_url": "amqp://guest:guest@localhost:5672/",
     "encryption_key": "thisis32byteencryptionkeyexample"
   }
   ```
5. Start service:
   ```sh
   go run main.go
   ```

## Run Producer with Docker

Create `config.json`:
```json
{
  "rabbitmq_url": "amqp://guest:guest@rabbitmq:5672/",
  "encryption_key": "thisis32byteencryptionkeyexample"
}
```
Start service:
```sh
docker network create rabbitmq_network

cd producer-service

docker-compose up -d --build
```
Verify:
```sh
docker ps
```
Stop service:
```sh
docker-compose down
```

## Troubleshooting
- **Port conflicts**: Ensure ports 5672, 15672 are available.
- **RabbitMQ connection issues**:
  ```sh
  docker logs rabbitmq-producer
  ```

# Consumer Service

Fetches data from RabbitMQ and stores it in PostgreSQL & Redis.

## Prerequisites
- **Docker & Docker Compose** ([Install](https://docs.docker.com/get-docker/))
- **Go** ([Download](https://go.dev/doc/install))
- **PostgreSQL** ([Download](https://www.postgresql.org/download/))
- **Redis** ([Download](https://redis.io/docs/getting-started/))
- **RabbitMQ** ([Download](https://www.rabbitmq.com/download.html))
- **cURL** (for API testing)

### Check Running Services
```sh
docker ps | grep -E "redis|rabbitmq|postgres"
```
Stop conflicts:
```sh
docker stop <container_id> && docker rm <container_id>
```

## Run Consumer Locally

1. Navigate to `consumer-service`:
```sh
cd consumer-service
```
2. Install dependencies:
```sh
go mod tidy
```
3. Ensure PostgreSQL, Redis, and RabbitMQ are running.
4. Create `config.json`:
```json
{
  "rabbitmq_url": "amqp://guest:guest@localhost:5672/",
  "redis_address": "localhost:6379",
  "postgresql_host": "localhost",
  "postgresql_port": "5432",
  "postgresql_user": "postgres",
  "postgresql_password": "password",
  "postgresql_dbname": "testdb",
  "encryption_key": "thisis32byteencryptionkeyexample"
}
```
5. Start service:
```sh
go run main.go
```

## Run Consumer with Docker

Create `config.json`:
```json
{
  "rabbitmq_url": "amqp://guest:guest@rabbitmq:5672/",
  "redis_address": "redis:6379",
  "postgresql_host": "postgres",
  "postgresql_port": "5432",
  "postgresql_user": "postgres",
  "postgresql_password": "password",
  "postgresql_dbname": "testdb",
  "encryption_key": "thisis32byteencryptionkeyexample"
}
```
Start service:
```sh
cd consumer-service

docker-compose up -d --build
```
Check running services:
```sh
docker ps
```
Stop service:
```sh
docker-compose down
```

## Troubleshooting
- **RabbitMQ issues**:
  ```sh
  docker logs rabbitmq-consumer
  ```
- **PostgreSQL issues**:
  ```sh
  docker logs postgres-consumer
  ```

# Database Setup

Create `users` table:
```sh
docker exec -it postgres psql -U postgres -d testdb
```
Run:
```sql
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  first_name VARCHAR(255),
  last_name VARCHAR(255),
  email_address VARCHAR(255) NOT NULL,
  created_at VARCHAR(255),
  merged_at VARCHAR(255),
  deleted_at VARCHAR(255),
  parent_user_id FLOAT
);
```
Exit:
```sh
\q
```

# API Usage

### Upload CSV (Producer API)
```sh
curl --location 'http://localhost:8080/upload-csv' --form 'file=@"/path/to/Users1.csv"'
```

### Retrieve User Data (Consumer API)
```sh
curl --location 'http://localhost:8081/users?first_name=Grace&last_name=Taylor'
```

🚀 Happy Coding!

