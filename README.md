# BigDataForge - Scalable RESTful API for Complex JSON Data Processing

## 📌 Overview
BigDataForge is a **scalable RESTful API** built with **Golang**, designed to **manage and process complex structured data** in **JSON format**. Adopting a **distributed systems approach**, it ensures **high availability and reliability** while integrating with technologies like **Redis, Elasticsearch, and RabbitMQ** for optimized data storage, indexing, and message queueing.

This project supports **comprehensive CRUD operations**, **advanced conditional updates**, and **search capabilities**, making it ideal for handling **hierarchical data relationships** using **parent-child indexing in Elasticsearch**.

---

## 🚀 Features
- **🌍 RESTful API**: Structured endpoints for managing JSON data.
- **📝 CRUD Operations**: Full support for **create, read, update, delete, and patch** operations with schema validation.
- **🔍 Elasticsearch Integration**: Efficient **parent-child indexing** for hierarchical data relationships.
- **📩 Queueing System**: Uses **RabbitMQ** to ensure reliable data ingestion.
- **✅ Data Validation**: Enforces **JSON schema validation** for data integrity.
- **🔐 Security**: Implements **OAuth 2.0 authentication with GCP**.
- **⚡ Distributed System Architecture**: Utilizes **Redis for key-value storage** and **Elasticsearch for advanced indexing and search**.

---

## 🛠️ Technologies Used
- **Golang**: Core backend programming language.
- **Gin Framework**: Lightweight web framework for REST API development.
- **Redis**: In-memory key-value store for fast data access.
- **gojsonschema**: Library for JSON schema validation.
- **Postman**: Used for testing API endpoints.
- **Elasticsearch**: High-performance search and analytics engine.
- **RabbitMQ**: Message broker for asynchronous processing.

---

## 🔄 Data Flow
1️⃣ **Generate OAuth token** via authentication workflow.
2️⃣ **Validate API requests** using the received ID token.
3️⃣ **Create JSON object** using the `POST` HTTP method.
4️⃣ **Validate JSON object** against the defined JSON schema.
5️⃣ **De-structure hierarchical JSON** while storing in **Redis**.
6️⃣ **Enqueue object in RabbitMQ** for indexing.
7️⃣ **Dequeue from RabbitMQ** and index the data in **Elasticsearch**.
8️⃣ **Implement search queries** using **Kibana Console** for retrieving indexed data.

---

## 📥 Installation
### **Prerequisites**
- **Go 1.17+** (Download from [Go official site](https://go.dev/dl/))
- **Redis** (Install from [Redis official site](https://redis.io/))
- **Elasticsearch** (Install from [Elasticsearch official site](https://www.elastic.co/))
- **RabbitMQ** (Install from [RabbitMQ official site](https://www.rabbitmq.com/))
- **Refer to `example.env`** for setting up required environment variables.

### **Run with Docker**
```sh
git clone https://github.com/akhilk2802/BigDataForge.git
cd BigDataForge
docker-compose up --build
```

### **Run Manually**
```sh
git clone https://github.com/akhilk2802/BigDataForge.git
cd BigDataForge

# Install dependencies
go mod tidy

# Start Redis
redis-server

# Start Elasticsearch
elasticsearch

# Start RabbitMQ
rabbitmq-server

# Run the service
go run cmd/api/main.go
go run cmd/listener/main.go
```

---

## 🔗 API Endpoints
### **📌 Create a New Plan**
```http
POST /v1/plan
```
- Creates a new plan from the request body.

### **📌 Update an Existing Plan**
```http
PUT /v1/plan/{id}
```
- Updates an existing plan by **ID**.
- Requires a valid **ETag** in the `If-Match` HTTP header.

### **📌 Patch an Existing Plan**
```http
PATCH /v1/plan/{id}
```
- Partially updates a plan by **ID**.
- Requires a valid **ETag** in the `If-Match` HTTP header.

### **📌 Fetch an Existing Plan**
```http
GET /v1/plan/{id}
```
- Retrieves a plan by **ID**.
- Supports **ETag-based caching** with `If-None-Match` HTTP header.

### **📌 Delete an Existing Plan**
```http
DELETE /v1/plan/{id}
```
- Deletes a plan by **ID**.
- Requires a valid **ETag** in the `If-Match` HTTP header.

---

## 📞 Contact
For inquiries or contributions:
- **GitHub:** [yourusername](https://github.com/yourusername)
- **Email:** your.email@example.com

🚀 **BigDataForge - Powering Scalable & Efficient JSON Data Processing!**

# BigDataForge - Scalable RESTful API for Complex JSON Data Processing

## 📌 Overview
BigDataForge is a **scalable RESTful API** built with **Golang**, designed to **manage and process complex structured data** in **JSON format**. Adopting a **distributed systems approach**, it ensures **high availability and reliability** while integrating with technologies like **Redis, Elasticsearch, and RabbitMQ** for optimized data storage, indexing, and message queueing.

This project supports **comprehensive CRUD operations**, **advanced conditional updates**, and **search capabilities**, making it ideal for handling **hierarchical data relationships** using **parent-child indexing in Elasticsearch**.

---

## 🚀 Features
- **🌍 RESTful API**: Structured endpoints for managing JSON data.
- **📝 CRUD Operations**: Full support for **create, read, update, delete, and patch** operations with schema validation.
- **🔍 Elasticsearch Integration**: Efficient **parent-child indexing** for hierarchical data relationships.
- **📩 Queueing System**: Uses **RabbitMQ** to ensure reliable data ingestion.
- **✅ Data Validation**: Enforces **JSON schema validation** for data integrity.
- **🔐 Security**: Implements **OAuth 2.0 authentication with GCP**.
- **⚡ Distributed System Architecture**: Utilizes **Redis for key-value storage** and **Elasticsearch for advanced indexing and search**.

---

## 🛠️ Technologies Used
- **Golang**: Core backend programming language.
- **Gin Framework**: Lightweight web framework for REST API development.
- **Redis**: In-memory key-value store for fast data access.
- **gojsonschema**: Library for JSON schema validation.
- **Postman**: Used for testing API endpoints.
- **Elasticsearch**: High-performance search and analytics engine.
- **RabbitMQ**: Message broker for asynchronous processing.

---

## 🔄 Data Flow
1️⃣ **Generate OAuth token** via authentication workflow.

2️⃣ **Validate API requests** using the received ID token.

3️⃣ **Create JSON object** using the `POST` HTTP method.

4️⃣ **Validate JSON object** against the defined JSON schema.

5️⃣ **De-structure hierarchical JSON** while storing in **Redis**.

6️⃣ **Enqueue object in RabbitMQ** for indexing.

7️⃣ **Dequeue from RabbitMQ** and index the data in **Elasticsearch**.

8️⃣ **Implement search queries** using **Kibana Console** for retrieving indexed data.

---

## 📥 Installation
### **Prerequisites**
- **Go 1.17+** (Download from [Go official site](https://go.dev/dl/))
- **Redis** (Install from [Redis official site](https://redis.io/))
- **Elasticsearch** (Install from [Elasticsearch official site](https://www.elastic.co/))
- **RabbitMQ** (Install from [RabbitMQ official site](https://www.rabbitmq.com/))
- **Refer to `example.env`** for setting up required environment variables.

---

## Folder Structure
```plaintext
.
├── README.md
├── cmd
│   ├── api
│   │   └── main.go
│   └── listener
│       └── main.go
├── docker-compose.yml
├── dockerfile
├── dump.rdb
├── example.env
├── go.mod
├── go.sum
└── internal
    ├── controllers
    │   └── plan_controller.go
    ├── elastic
    │   └── factory.go
    ├── middlewares
    │   └── authentication.go
    ├── models
    │   └── plan.go
    ├── rabbitmq
    │   └── factory.go
    ├── routes
    │   └── api_routes.go
    ├── schemas
    │   ├── patch_plan_schema.json
    │   └── plan_schema.json
    ├── services
    │   └── plan_service.go
    ├── storage
    │   └── redis.go
    └── validators
        └── plan_validator.go
```

---

### **Run with Docker**
```sh
git clone https://github.com/akhilk2802/BigDataForge.git
cd BigDataForge
docker-compose up --build
```

### **Run Manually**
```sh
git clone https://github.com/akhilk2802/BigDataForge.git
cd BigDataForge

# Install dependencies
go mod tidy

# Start Redis
redis-server

# Start Elasticsearch
elasticsearch

# Start RabbitMQ
rabbitmq-server

# Run the service
go run cmd/api/main.go
go run cmd/listener/main.go
```

---

## 🔗 API Endpoints
### **📌 Create a New Plan**
```http
POST /v1/plan
```
- Creates a new plan from the request body.

### **📌 Update an Existing Plan**
```http
PUT /v1/plan/{id}
```
- Updates an existing plan by **ID**.
- Requires a valid **ETag** in the `If-Match` HTTP header.

### **📌 Patch an Existing Plan**
```http
PATCH /v1/plan/{id}
```
- Partially updates a plan by **ID**.
- Requires a valid **ETag** in the `If-Match` HTTP header.

### **📌 Fetch an Existing Plan**
```http
GET /v1/plan/{id}
```
- Retrieves a plan by **ID**.
- Supports **ETag-based caching** with `If-None-Match` HTTP header.

### **📌 Delete an Existing Plan**
```http
DELETE /v1/plan/{id}
```
- Deletes a plan by **ID**.
- Requires a valid **ETag** in the `If-Match` HTTP header.

---

🚀 **BigDataForge - Powering Scalable & Efficient JSON Data Processing!**

