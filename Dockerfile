# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o customer-service cmd/main.go

# Run stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/customer-service .
COPY .env .
COPY ./migrations /app/migrations

# Set the migrations directory for container environment
ENV MIGRATIONS_DIR=/app/migrations

EXPOSE 8080

CMD ["./customer-service"]
