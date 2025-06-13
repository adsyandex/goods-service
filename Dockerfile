FROM golang:1.18-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o goods-service .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/goods-service .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./goods-service"]