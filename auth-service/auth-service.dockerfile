# Stage 1: Build Go binary
FROM golang:1.25-alpine AS builder

RUN mkdir /app
WORKDIR /app
COPY . /app

# Set Go build parameters to target Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o auth-service ./cmd/api

# Stage 2: Create minimal runtime image
FROM alpine:latest

# Copy the Go binary from the builder stage
COPY --from=builder /app/auth-service /app/auth-service

# Make the binary executable
RUN chmod +x /app/auth-service

CMD ["/app/auth-service"]

LABEL authors="ngtronghao"
