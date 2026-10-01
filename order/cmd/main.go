package main

import (
	"log"
	orderpb "marketplace-api/gen/order"
	"net"
	"net/http"

	ordergrpc "order/internal/grpc"

	"google.golang.org/grpc"

	"order/internal/client"
	"order/internal/config"
	"order/internal/database"
	"order/internal/handler"
	"order/internal/repository"
	"order/internal/service"
)

func main() {
	cfg := config.Load()

	db, err := database.New(cfg.DatabaseDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	orderRepository := repository.NewOrderRepository(db)

	analyticsClient, err := client.NewAnalyticsClient(cfg.AnalyticsAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer analyticsClient.Close()

	productClient, err := client.NewProductClient(cfg.ProductAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer productClient.Close()

	orderService := service.NewOrderService(
		orderRepository,
		analyticsClient,
		productClient,
	)

	orderGRPCServer := ordergrpc.NewServer(orderService)

	grpcServer := grpc.NewServer()

	orderpb.RegisterOrderServiceServer(
		grpcServer,
		orderGRPCServer,
	)

	listener, err := net.Listen("tcp", ":9091")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		log.Printf("order grpc service started on :9091")

		if err := grpcServer.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	orderHandler := handler.NewOrderHandler(orderService)

	mux := http.NewServeMux()

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	mux.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			orderHandler.Create(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/orders/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			orderHandler.GetByID(w, r)
		case http.MethodPatch:
			orderHandler.Cancel(w, r)
		case http.MethodDelete:
			orderHandler.Delete(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Printf("order service started on %s", cfg.ServerAddress)

	if err := http.ListenAndServe(cfg.ServerAddress, mux); err != nil {
		log.Fatal(err)
	}
}
