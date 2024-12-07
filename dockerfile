# Use the official Go image
FROM golang:latest as builder

# Set the working directory
WORKDIR /app

# Copy Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the application source code
COPY . .

# Build both applications
RUN go build -o api ./cmd/api/main.go
RUN go build -o listener ./cmd/listener/main.go

# Use a minimal image for deployment
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/api .
COPY --from=builder /app/listener .

# Expose the ports
EXPOSE 8080

# Set environment variables
ENV REDIS_ADDR=redis:6379 \
    REDIS_PASSWORD= \
    REDIS_DB=0 \
    GOOGLE_CLIENT_ID=218327216067-1rqvp86sef5d5dcr6tcu8mfblvodmpgt.apps.googleusercontent.com \
    ELASTICSEARCH_URL=http://elasticsearch:9200 \
    RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/

# Use a process manager for multiple services
CMD ["sh", "-c", "./api & ./listener"]