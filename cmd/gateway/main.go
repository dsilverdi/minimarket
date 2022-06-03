package main

import (
	"log"
	"minimarket/gateway"
	"minimarket/gateway/api"
	"minimarket/gateway/client"

	"google.golang.org/grpc"
)

func main() {
	var user_conn, product_conn *grpc.ClientConn
	var err error

	user_conn, err = grpc.Dial("usermodule:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer user_conn.Close()

	product_conn, err = grpc.Dial("productmodule:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer user_conn.Close()

	usercl := client.NewUserClient(user_conn)
	productcl := client.NewProductClient(product_conn)
	svc := gateway.New(usercl, productcl)
	e := api.NewHttpAPI(svc)

	e.Logger.Fatal(e.Start(":8000"))
}
