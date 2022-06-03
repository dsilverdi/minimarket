package mocks

import (
	"context"
	"minimarket/gateway"
	"minimarket/pkg/errors"
	"minimarket/pkg/rest"
	"minimarket/product/api"
)

type ProductClient struct {
	productlist []*api.ProductData
	commentlist []*api.CommentData
}

func NewProductClient(products []*api.ProductData, comments []*api.CommentData) gateway.ProductClientInterface {
	return &ProductClient{
		productlist: products,
		commentlist: comments,
	}
}

func (cl *ProductClient) GetProducts(ctx context.Context, name string, category string, query rest.QueryParameter) (*api.ProductResponse, error) {
	return &api.ProductResponse{
		Products: cl.productlist,
		Comments: cl.commentlist,
	}, nil
}

func (cl *ProductClient) WriteComment(ctx context.Context, productid int, replyid int, owner string, message string) (*api.CommentData, error) {
	if productid == 0 {
		return nil, errors.ErrCreateEntity
	}
	return nil, nil
}
