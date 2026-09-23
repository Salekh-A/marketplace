FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o order ./cmd

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/order .

EXPOSE 8081
EXPOSE 9091

CMD ["./order"]