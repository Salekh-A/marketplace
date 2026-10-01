FROM golang:1.25 AS builder

WORKDIR /app

COPY order/go.mod order/go.sum ./order/
COPY api/grpc/go.mod api/grpc/go.sum ./api/grpc/

RUN cd order && go mod download

COPY order ./order
COPY api/grpc ./api/grpc

WORKDIR /app/order

RUN go build -o order ./cmd

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/order/order .

EXPOSE 8081
EXPOSE 9091

CMD ["./order"]