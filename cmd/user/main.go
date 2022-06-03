package main

import (
	"log"
	"minimarket/pkg/uuid"
	"minimarket/user"
	"minimarket/user/api"
	"minimarket/user/database"
	"net"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
)

func main() {
	db, err := database.Connect(database.Config{
		Host:        "usermoduledb",
		Port:        "3306",
		User:        "root",
		Pass:        "usermodule",
		Name:        "usermodule",
		SSLMode:     "disable",
		SSLCert:     "",
		SSLKey:      "",
		SSLRootCert: "",
	})
	if err != nil {
		log.Fatal("Error Connecting DB", err)
	}
	defer db.Close()

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	//setup grpc server
	grpcServer := grpc.NewServer()
	UserAPI := NewService(db)

	api.RegisterAuthServiceServer(grpcServer, UserAPI)

	errs := make(chan error)

	log.Print("Starting User Module")
	errs <- grpcServer.Serve(lis)

	log.Fatal("exit", <-errs)
}

func NewService(db *sqlx.DB) api.AuthServiceServer {
	User := database.NewUsersRepository(db)
	IDprov := uuid.New()
	return user.NewServer(User, IDprov)
}
