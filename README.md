# BigDataForge

BigDataForge is a RESTful API project designed to manage and process complex structured data in JSON format. Built with Golang, it leverages a distributed systems approach to ensure scalability and reliability, integrating technologies like Redis, Elasticsearch, and RabbitMQ for optimized data storage, indexing, and message queueing.

This project provides comprehensive support for CRUD operations, advanced conditional updates, and search capabilities. It demonstrates parent-child indexing and search using Elasticsearch, making it ideal for scenarios requiring hierarchical data relationships and efficient query execution.


## Features

- **RESTful APIs:** Flexible endpoints to handle structured JSON data.
- **CRUD Operations:** Create, read, update, delete, and patch capabilities with advanced validation.
- **Elasticsearch Integration:** Parent-child indexing and search for efficient hierarchical data retrieval.
- **Queueing System:** RabbitMQ for handling message queues and ensuring reliable data ingestion.
- **Data Validation:** Schema validation to ensure data integrity.
- **Security:** Authentication using GCP OAuth2.0
- **Distributed System Architecture:** Leveraging Redis for key-value storage and Elasticsearch for advanced indexing and search.

## Technologies Used

- **Golang**: Core programming language.
- **Gin Framework**: Lightweight web framework for building REST APIs.
- **Redis**: In-memory key-value database used for storing plan data.
- **gojsonschema**: JSON Schema validation library.
- **Postman**: Used for testing API endpoints.


## Data Flow
1. Generate OAuth token using authorization work flow
2. Validate further API requests using the received ID token
3. Create JSON Object using the `POST` HTTP method
4. Validate incoming JSON Object using the respective JSON Schema
5. De-Structure hierarchial JSON Object while storing in Redis key-value store
6. Enqueue object in RabbitMQ queue to index the object
7. Dequeue from RabbitMQ queue and index data in ElasticServer
8. Implement Search queries using Kibana Console to retrieve indexed data


## Installation

### Prerequisites

- **Go 1.17 or later**: Ensure that Golang is installed on your system. You can download it from [here](https://golang.org/dl/).
- **Redis**: Make sure Redis is installed and running. Follow the installation instructions from the [Redis website](https://redis.io/download).


**Refer to example.env for setting up the required env variables** 

### Steps: to run with `Docker`

1. Clone the repository:
   ```bash
   git clone https://github.com/akhilk2802/BigDataForge.git
   cd BigDataForge
   ```
2. Build and run the docker-compose
   ```bash
   docker-compose up --build
   ```


### Steps: to run `Manually`

1. Clone the repository:

   ```bash
   git clone https://github.com/akhilk2802/BigDataForge.git
   cd BigDataForge
   ```
2. Install the dependencies:
    
    ```bash
    go mod tidy
    ```

3.  Start the redis server:
    
    ```bash
    redis-server
    ```

4. Start the elastic search:

    ```bash
    elasticsearch
    ```

5. Start the RabbitMQ:

    ```bash
    rabbitmq-server
    ```

6.  Run the service:
    
    ```bash
    go run cmd/api/main.go
    go run cmd/listener/main.go
    ```

## API Endpoints


1. POST `/v1/plan` - Creates a new plan provided in the request body
2. PUT `/v1/plan/{id}` - Updates an existing plan provided by the id
    - A valid Etag for the object should also be provided in the `If-Match` HTTP Request Header
3. PATCH `/v1/plan/{id}` - Patches an existing plan provided by the id
    - A valid Etag for the object should also be provided in the `If-Match` HTTP Request Header
4. GET `/v1/plan/{id}` - Fetches an existing plan provided by the id
    - An Etag for the object can be provided in the `If-None-Match` HTTP Request Header
    - If the request is successful, a valid Etag for the object is returned in the `ETag` HTTP Response Header
5. DELETE `/v1/plan/{id}` - Deletes an existing plan provided by the id
    - A valid Etag for the object should also be provided in the `If-Match` HTTP Request Header