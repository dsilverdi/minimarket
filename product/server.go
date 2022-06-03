package product

import (
	"context"
	"log"
	"minimarket/pkg/errors"
	"minimarket/product/api"
	"time"
)

type ProductServer struct {
	Product ProductRepository
	api.UnimplementedProductServiceServer
}

func NewServer(product ProductRepository) api.ProductServiceServer {
	return &ProductServer{
		Product: product,
	}
}

func (srv *ProductServer) GetProducts(ctx context.Context, req *api.ProductRequest) (*api.ProductResponse, error) {
	queries, err := srv.Product.Read(ctx, req.Name, req.Category, int(req.Limit), int(req.Page))
	if err != nil {
		log.Print(err)
		return nil, errors.ErrViewEntity
	}

	prodidMap := make(map[int32]int)
	var productDatas []*api.ProductData
	var commentDatas []*api.CommentData
	for _, q := range queries {
		if _, ok := prodidMap[q.ProductID]; !ok {
			prodidMap[q.ProductID] = 1
			productDatas = append(productDatas, &api.ProductData{
				Id:          int32(q.ProductID),
				Name:        q.ProductName,
				Description: q.ProductDescription,
				Category:    q.ProductCategory,
				Picture:     q.ProductPicture,
			})
		}

		if q.CommentID != nil {
			commentDatas = append(commentDatas, &api.CommentData{
				Id:        *q.CommentID,
				Productid: q.ProductID,
				Message:   *q.CommentMessage,
				Owner:     *q.CommentAuthor,
				Replyid:   *q.ReplyID,
				Createdat: q.CreatedAt.Format(time.RFC3339),
			})
		}

	}

	return &api.ProductResponse{
		Products: productDatas,
		Comments: commentDatas,
	}, nil
}

func (srv *ProductServer) WriteComment(ctx context.Context, req *api.WriteCommentRequest) (*api.CommentData, error) {
	err := srv.Product.WriteComment(ctx, int(req.Productid), int(req.Replyid), req.Message, req.Owner)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return &api.CommentData{
		Productid: req.Productid,
		Message:   req.Message,
		Owner:     req.Owner,
	}, nil
}
