# BigDataForge

BigDataForge is a REST API service built using Golang's Gin framework and Redis as a key-value store. This service is designed to handle any structured data in JSON format and provides CRUD operations along with JSON Schema validation for incoming data.

## Features

- **POST /plans**: Create a new plan with structured data.
- **GET /plans**: Retrieve a plan by its ID.
- **DELETE /plans**: Delete a plan by its ID.
- **GET /plans/conditional**: Retrieve a plan conditionally based on attributes like `planType`.
- **Validation**: JSON schema validation for incoming requests.
- **Key-Value Store**: Redis is used for storing plan data as key-value pairs.

## Technologies Used

- **Golang**: Core programming language.
- **Gin Framework**: Lightweight web framework for building REST APIs.
- **Redis**: In-memory key-value database used for storing plan data.
- **gojsonschema**: JSON Schema validation library.
- **Postman**: Used for testing API endpoints.

## Installation

### Prerequisites

- **Go 1.17 or later**: Ensure that Golang is installed on your system. You can download it from [here](https://golang.org/dl/).
- **Redis**: Make sure Redis is installed and running. Follow the installation instructions from the [Redis website](https://redis.io/download).

### Steps

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

4.  Run the service:
    
    ```bash
    go run cmd/api/main.go
    ```

## API Endpoints

1. Create a New Plan (POST/plans)
2. Retrieve a Plan by ID (GET /plans?id=<planId>)
3. Delete a Plan by ID (DELETE /plans?id=<planId>)
4. Delete a Plan by ID (DELETE /plans/conditional/?id=<planId>)