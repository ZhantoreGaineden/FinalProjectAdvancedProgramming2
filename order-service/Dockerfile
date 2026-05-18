FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN mkdir -p /app/bin && go build -o /app/bin/order-service ./order-service/cmd

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/bin/order-service .

EXPOSE 50053

CMD ["./order-service"]
