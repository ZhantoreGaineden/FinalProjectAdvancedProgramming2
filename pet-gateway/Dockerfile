FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN mkdir -p /app/bin && go build -o /app/bin/pet-gateway ./pet-gateway/cmd

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/bin/pet-gateway .

EXPOSE 8081

CMD ["./pet-gateway"]
