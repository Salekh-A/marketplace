FROM golang:1.25 AS builder

WORKDIR /app

COPY analytics/go.mod analytics/go.sum ./analytics/
COPY api/grpc/go.mod api/grpc/go.sum ./api/grpc/

RUN cd analytics && go mod download

COPY analytics ./analytics
COPY api/grpc ./api/grpc

WORKDIR /app/analytics

RUN go build -o analytics ./cmd

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/analytics/analytics .

EXPOSE 8083
EXPOSE 9093

CMD ["./analytics"]