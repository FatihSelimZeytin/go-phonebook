# Stage 1: Build the Go app
FROM golang:1.22-alpine AS builder

# Install git if needed
RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o go-phonebook .

# Stage 2: Final tiny image
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/go-phonebook .

ENTRYPOINT ["/app/go-phonebook"]
