# Stage 1: Build Go binary
FROM golang:1.25-alpine AS builder

RUN mkdir /app
WORKDIR /app
COPY . /app

# Set Go build parameters to target Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o broker-service ./cmd/api

# Stage 2: Create minimal runtime image
FROM alpine:latest

# Copy the Go binary from the builder stage
COPY --from=builder /app/broker-service /app/broker-service

# Make the binary executable
RUN chmod +x /app/broker-service

CMD ["/app/broker-service"]

LABEL authors="ngtronghao"
