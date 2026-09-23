package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"analytics/internal/config"
	"analytics/internal/database"
	analyticsgrpc "analytics/internal/grpc"
	"analytics/internal/handler"
	"analytics/internal/repository"
	"analytics/internal/service"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()

	db, err := database.New(cfg.DatabaseDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddress,
	})

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatal(err)
	}
	defer redisClient.Close()

	analyticsRepository := repository.NewAnalyticsRepository(db)

	analyticsService := service.NewAnalyticsService(
		analyticsRepository,
		redisClient,
	)

	analyticsHandler := handler.NewAnalyticsHandler(
		analyticsService,
	)

	analyticsGRPCServer := analyticsgrpc.NewServer(
		analyticsService,
	)

	grpcServer := grpc.NewServer()

	analyticsgrpc.RegisterAnalyticsServiceServer(
		grpcServer,
		analyticsGRPCServer,
	)

	listener, err := net.Listen("tcp", ":9093")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		log.Printf("analytics grpc service started on %s", ":9093")

		if err := grpcServer.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	mux := http.NewServeMux()

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	mux.HandleFunc("/analytics/revenue", analyticsHandler.Revenue)
	mux.HandleFunc("/analytics/orders", analyticsHandler.Orders)
	mux.HandleFunc("/analytics/average-check", analyticsHandler.AverageCheck)
	mux.HandleFunc("/analytics/top-products", analyticsHandler.TopProducts)

	log.Printf(
		"analytics service started on %s",
		cfg.ServerAddress,
	)

	if err := http.ListenAndServe(
		cfg.ServerAddress,
		mux,
	); err != nil {
		log.Fatal(err)
	}
}
