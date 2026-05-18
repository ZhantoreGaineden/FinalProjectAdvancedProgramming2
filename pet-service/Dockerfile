FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN mkdir -p /app/bin && go build -o /app/bin/pet-service ./pet-service/cmd

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/bin/pet-service .

EXPOSE 50051

CMD ["./pet-service"]
