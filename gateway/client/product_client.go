package client

import (
	"context"
	"log"
	"minimarket/gateway"
	"minimarket/pkg/rest"
	"minimarket/product/api"

	"google.golang.org/grpc"
)

type ProductClient struct {
	client api.ProductServiceClient
}

func NewProductClient(conn *grpc.ClientConn) gateway.ProductClientInterface {
	return &ProductClient{
		client: api.NewProductServiceClient(conn),
	}
}

func (cl *ProductClient) GetProducts(ctx context.Context, name string, category string, query rest.QueryParameter) (*api.ProductResponse, error) {

	payload := &api.ProductRequest{
		Name:     name,
		Category: category,
		Limit:    int32(query.Limit),
		Page:     int32(query.Page),
	}

	res, err := cl.client.GetProducts(ctx, payload)
	if err != nil {
		log.Print("client error ", err)
		return nil, err
	}

	return res, nil
}

func (cl *ProductClient) WriteComment(ctx context.Context, productid int, replyid int, owner string, message string) (*api.CommentData, error) {
	payload := &api.WriteCommentRequest{
		Productid: int32(productid),
		Replyid:   int32(replyid),
		Owner:     owner,
		Message:   message,
	}

	res, err := cl.client.WriteComment(ctx, payload)
	if err != nil {
		log.Print("client error ", err)
		return nil, err
	}

	return res, nil
}
