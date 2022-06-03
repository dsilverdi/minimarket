package client

import (
	"context"
	"log"
	"minimarket/gateway"
	"minimarket/user/api"

	"google.golang.org/grpc"
)

type UserClient struct {
	client api.AuthServiceClient
}

func NewUserClient(conn *grpc.ClientConn) gateway.UserClientInterface {
	return &UserClient{
		client: api.NewAuthServiceClient(conn),
	}
}

func (cl *UserClient) RegisterUser(ctx context.Context, email string, password string) (*api.UserIdentity, error) {
	payload := api.UserRequest{
		Email:    email,
		Password: password,
	}

	res, err := cl.client.Register(ctx, &payload)
	if err != nil {
		log.Print("client error ", err)
		return nil, err
	}

	return res, nil
}

func (cl *UserClient) AuthorizeUser(ctx context.Context, email string, password string) (*api.AuthResponse, error) {
	payload := api.UserRequest{
		Email:    email,
		Password: password,
	}

	res, err := cl.client.Authorize(ctx, &payload)
	if err != nil {
		log.Print("client error ", err)
		return nil, err
	}

	return res, nil
}
