FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN mkdir -p /app/bin && go build -o /app/bin/user-gateway ./user-gateway/cmd

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/bin/user-gateway .

EXPOSE 8082

CMD ["./user-gateway"]
