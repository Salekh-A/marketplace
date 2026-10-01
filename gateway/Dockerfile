FROM golang:1.25 AS builder

WORKDIR /app

COPY gateway/go.mod gateway/go.sum ./gateway/
COPY api/grpc/go.mod api/grpc/go.sum ./api/grpc/

RUN cd gateway && go mod download

COPY gateway ./gateway
COPY api/grpc ./api/grpc

WORKDIR /app/gateway

RUN go build -o gateway ./cmd

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/gateway/gateway .

EXPOSE 8080

CMD ["./gateway"]