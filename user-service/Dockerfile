FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN mkdir -p /app/bin && go build -o /app/bin/user-service ./user-service/cmd

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/bin/user-service .

EXPOSE 50052

CMD ["./user-service"]
