FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN mkdir -p /app/bin && go build -o /app/bin/order-gateway ./order-gateway/cmd

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/bin/order-gateway .

EXPOSE 8083

CMD ["./order-gateway"]
