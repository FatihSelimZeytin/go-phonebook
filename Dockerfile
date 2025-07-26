# Stage 1: Build the Go app
FROM golang:1.23 AS builder

# Install git if needed
#RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o go-phonebook .
RUN ls -all

# Stage 2: Final tiny image
FROM golang:1.22-alpine

WORKDIR /app
COPY --from=builder /app/go-phonebook .
COPY .env .env

EXPOSE 8080

CMD ["./go-phonebook"]
