# Use the official Go image for building
FROM golang:latest as builder

# Set the working directory
WORKDIR /app

# Copy Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the application source code
COPY . .

# Build both services
RUN go build -o api ./cmd/api/main.go
RUN go build -o listener ./cmd/listener/main.go

# Use a minimal image for deployment
FROM alpine:latest
WORKDIR /app

# Copy built binaries
COPY --from=builder /app/api .
COPY --from=builder /app/listener .

# Expose the ports
EXPOSE 8080

# Set environment variables
ENV REDIS_ADDR=redis:6379 \
    REDIS_PASSWORD= \
    REDIS_DB=0 \
    GOOGLE_CLIENT_ID=xxxx \
    ELASTICSEARCH_URL=http://elasticsearch:9200 \
    RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/

# Command placeholder - overridden by docker-compose.yml
CMD ["./api"]