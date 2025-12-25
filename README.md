# kafka-go
# ☕ Coffee Order System (Kafka + Go)

A simple event-driven coffee ordering system built with **Go** and **Apache Kafka**.  
The project demonstrates how to use Kafka producers and consumers using the **Sarama** library.

## 📌 Overview

This project consists of two Go applications:

1.  **Order Service (Producer)**
    
    *   Exposes an HTTP API to place coffee orders
        
    *   Publishes orders to a Kafka topic (`coffee_order`)
        
2.  **Coffee Worker (Consumer)**
    
    *   Consumes messages from the Kafka topic
        
    *   Simulates brewing coffee by processing incoming orders
        

* * *

## 🏗 Architecture

`Client   |   | POST /order   v Order Service (HTTP API)   |   | Kafka Producer   v Kafka Topic: coffee_order   |   | Kafka Consumer   v Coffee Worker (Brews Coffee ☕)`

* * *

## ⚙️ Tech Stack

*   **Go (Golang)**
    
*   **Apache Kafka**
    
*   **Sarama** (Kafka client for Go)
    
*   **net/http** (HTTP server)
    

* * *

## 📂 Project Structure

`. ├── producer/ │   └── main.go        # HTTP API + Kafka Producer ├── consumer/ │   └── main.go        # Kafka Consumer (Worker) └── README.md`

* * *

## 🔧 Prerequisites

Make sure you have the following installed:

*   Go 1.20+
    
*   Apache Kafka
    
*   Zookeeper (if using Kafka < 3.x)
    
*   Docker (optional, recommended)
    

### Kafka (Docker example)

`docker run -d --name zookeeper -p 2181:2181 zookeeper docker run -d --name kafka -p 9092:9092 \   -e KAFKA_ZOOKEEPER_CONNECT=host.docker.internal:2181 \   -e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://localhost:9092 \   confluentinc/cp-kafka`

* * *

## 📦 Kafka Topic

Create the topic before running the apps:

`kafka-topics.sh --create \   --topic coffee_order \   --bootstrap-server localhost:9092 \   --partitions 1 \   --replication-factor 1`

* * *

## ▶️ Running the Applications

### 1️⃣ Start the Kafka Consumer (Worker)

`go run consumer/main.go`

Output example:

`Your server is running. Press Ctrl+C to exit. 1 Brewing coffee {"consumer_name":"Alice","coffee_type":"Latte"}`

* * *

### 2️⃣ Start the Order Service (Producer)

`go run producer/main.go`

Server runs on:

`http://localhost:8000`

* * *

## 📡 API Usage

### Place a Coffee Order

**Endpoint**

`POST /order`

**Request Body**

`{   "consumer_name": "Alice",   "coffee_type": "Latte" }`

**cURL Example**

`curl -X POST http://localhost:8000/order \   -H "Content-Type: application/json" \   -d '{"consumer_name":"Alice","coffee_type":"Latte"}'`

**Response**

`{   "successs": true,   "msg": "Order for Customer:Alice send successfully" }`

* * *

## 🧠 How It Works

### Producer

*   Accepts HTTP POST requests
    
*   Serializes orders to JSON
    
*   Publishes messages to Kafka using a **SyncProducer**
    

### Consumer

*   Subscribes to partition `0` of `coffee_order`
    
*   Reads messages from the oldest offset
    
*   Gracefully shuts down on `Ctrl+C`
    
*   Logs and processes each order
    

* * *

## 🛑 Graceful Shutdown

The consumer listens for:

*   `SIGINT`
    
*   `SIGTERM`
    

On shutdown:

*   Stops consuming messages
    
*   Prints total processed message count
    
*   Closes Kafka connections safely
    

* * *

## 🚀 Future Improvements

*   Use **Consumer Groups**
    
*   Add authentication & validation
    
*   Support multiple partitions
    
*   Add retries & dead-letter queue
    
*   Docker Compose setup
    
*   Structured logging
    

* * *

## 📜 License

This project is open-source and available under the MIT License.
