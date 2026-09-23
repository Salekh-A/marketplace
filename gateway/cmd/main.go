package main

import (
	"log"
	"net/http"

	"gateway/internal/config"
	"gateway/internal/grpc"
	"gateway/internal/handler"
)

func main() {
	cfg := config.Load()

	orderClient, err := grpc.NewOrderClient(
		cfg.OrderAddress,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer orderClient.Close()

	productClient, err := grpc.NewProductClient(
		cfg.ProductAddress,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer productClient.Close()

	analyticsClient, err := grpc.NewAnalyticsClient(
		cfg.AnalyticsAddress,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer analyticsClient.Close()

	gatewayHandler := handler.NewHandler(
		orderClient,
		productClient,
		analyticsClient,
	)

	mux := http.NewServeMux()

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	mux.HandleFunc("/orders", gatewayHandler.CreateOrder)
	mux.HandleFunc("/orders/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			gatewayHandler.CancelOrder(w, r)
			return
		}

		gatewayHandler.GetOrder(w, r)
	})

	mux.HandleFunc("/products", gatewayHandler.GetProducts)
	mux.HandleFunc("/products/", gatewayHandler.GetProduct)

	mux.HandleFunc(
		"/analytics/revenue",
		gatewayHandler.GetRevenue,
	)

	mux.HandleFunc(
		"/analytics/orders",
		gatewayHandler.GetOrders,
	)

	log.Printf(
		"gateway started on %s",
		cfg.ServerAddress,
	)

	if err := http.ListenAndServe(
		cfg.ServerAddress,
		mux,
	); err != nil {
		log.Fatal(err)
	}
}
