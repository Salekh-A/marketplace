package main

import (
	"log"
	"net"
	"net/http"

	"product/internal/config"
	"product/internal/database"
	grpcserver "product/internal/grpc"
	"product/internal/handler"
	"product/internal/repository"
	"product/internal/service"

	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()

	db, err := database.New(cfg.DatabaseDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	productRepository := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepository)
	productHandler := handler.NewProductHandler(productService)

	productGRPCServer := grpcserver.NewServer(productService)

	grpcServer := grpc.NewServer()

	grpcserver.RegisterProductServiceServer(
		grpcServer,
		productGRPCServer,
	)

	listener, err := net.Listen("tcp", ":9092")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		log.Printf("product grpc service started on %s", ":9092")

		if err := grpcServer.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	mux := http.NewServeMux()

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	mux.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			productHandler.Create(w, r)
		case http.MethodGet:
			productHandler.GetAll(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/products/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			productHandler.GetByID(w, r)
		case http.MethodPatch:
			productHandler.Update(w, r)
		case http.MethodDelete:
			productHandler.Delete(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Printf("product service started on %s", cfg.ServerAddress)

	if err := http.ListenAndServe(cfg.ServerAddress, mux); err != nil {
		log.Fatal(err)
	}
}
