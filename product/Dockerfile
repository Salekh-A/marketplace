FROM golang:1.25 AS builder

WORKDIR /app

COPY product/go.mod product/go.sum ./product/
COPY api/grpc/go.mod api/grpc/go.sum ./api/grpc/

RUN cd product && go mod download

COPY product ./product
COPY api/grpc ./api/grpc

WORKDIR /app/product

RUN go build -o product ./cmd

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/product/product .

EXPOSE 8082
EXPOSE 9092

CMD ["./product"]