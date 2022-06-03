package main

import (
	"log"
	"minimarket/product"
	"minimarket/product/api"
	"minimarket/product/database"
	"net"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
)

func main() {
	db, err := database.Connect(database.Config{
		Host:        "productmoduledb",
		Port:        "3306",
		User:        "root",
		Pass:        "productmodule",
		Name:        "productmodule",
		SSLMode:     "disable",
		SSLCert:     "",
		SSLKey:      "",
		SSLRootCert: "",
	})
	if err != nil {
		log.Fatal("Error Connecting DB", err)
	}
	defer db.Close()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	//initialize service
	grpcServer := grpc.NewServer()
	productAPI := NewService(db)

	api.RegisterProductServiceServer(grpcServer, productAPI)

	errs := make(chan error)

	log.Print("Starting Product Module")
	errs <- grpcServer.Serve(lis)

	log.Fatal("exit", <-errs)
}

func NewService(db *sqlx.DB) api.ProductServiceServer {
	Product := database.NewProductRepository(db)
	return product.NewServer(Product)
}
