package product_test

import (
	"context"
	"fmt"
	"minimarket/pkg/errors"
	"minimarket/product"
	"minimarket/product/api"
	"minimarket/product/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	existproduct1  = "baju abc"
	existproduct2  = "celana abc"
	existproduct3  = "iphone 12"
	existcategory1 = "pakaian"
	existcategory2 = "elektronik"
	currentTime    = time.Now()
	author         = "user"
	Comment        = []int32{1, 2, 3, 4, 5, 6, 7}
	ReplyIDNUL     = int32(0)
	ReplyID1       = int32(1)
	ReplyID6       = int32(6)
	msg            = "msg"
	defaultLimit   = 15
)

func NewService() api.ProductServiceServer {
	productmocks := mocks.NewProductRepository(
		[]product.ProductQuery{
			{ProductID: 1, ProductName: existproduct1, ProductCategory: existcategory1, ProductDescription: "sample", ProductPicture: "url", CommentID: &Comment[0], ReplyID: &ReplyIDNUL, CommentMessage: &msg, CommentAuthor: &author, CreatedAt: &currentTime},
			{ProductID: 1, ProductName: existproduct1, ProductCategory: existcategory1, ProductDescription: "sample", ProductPicture: "url", CommentID: &Comment[1], ReplyID: &ReplyID1, CommentMessage: &msg, CommentAuthor: &author, CreatedAt: &currentTime},
			{ProductID: 1, ProductName: existproduct1, ProductCategory: existcategory1, ProductDescription: "sample", ProductPicture: "url", CommentID: &Comment[2], ReplyID: &ReplyID1, CommentMessage: &msg, CommentAuthor: &author, CreatedAt: &currentTime},
			{ProductID: 2, ProductName: existproduct2, ProductCategory: existcategory1, ProductDescription: "sample", ProductPicture: "url", CommentID: &Comment[3], ReplyID: &ReplyIDNUL, CommentMessage: &msg, CommentAuthor: &author, CreatedAt: &currentTime},
			{ProductID: 2, ProductName: existproduct2, ProductCategory: existcategory1, ProductDescription: "sample", ProductPicture: "url", CommentID: &Comment[4], ReplyID: &ReplyIDNUL, CommentMessage: &msg, CommentAuthor: &author, CreatedAt: &currentTime},
			{ProductID: 3, ProductName: existproduct3, ProductCategory: existcategory2, ProductDescription: "sample", ProductPicture: "url", CommentID: &Comment[5], ReplyID: &ReplyIDNUL, CommentMessage: &msg, CommentAuthor: &author, CreatedAt: &currentTime},
			{ProductID: 3, ProductName: existproduct3, ProductCategory: existcategory2, ProductDescription: "sample", ProductPicture: "url", CommentID: &Comment[6], ReplyID: &ReplyID6, CommentMessage: &msg, CommentAuthor: &author, CreatedAt: &currentTime},
		},
	)

	return product.NewServer(productmocks)
}

func TestGetProducts(t *testing.T) {
	svc := NewService()
	ctx := context.Background()

	cases := []struct {
		desc     string
		name     string
		category string
		limit    int32
		err      error
	}{
		{"get all", "", "", int32(defaultLimit), nil},
		{"get specific name", existproduct1, "", int32(defaultLimit), nil},
		{"get specific category", "", existcategory2, int32(defaultLimit), nil},
		{"get with limit", "", "", 2, nil},
	}

	for _, tc := range cases {
		payload := api.ProductRequest{
			Name:     tc.name,
			Category: tc.category,
			Limit:    tc.limit,
			Page:     1,
		}
		_, err := svc.GetProducts(ctx, &payload)
		assert.Equal(t, err, tc.err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
	}
}

func TestWriteComments(t *testing.T) {
	svc := NewService()
	ctx := context.Background()

	cases := []struct {
		desc    string
		prodid  int
		replyid int
		message string
		err     error
	}{
		{"write comments", 1, 0, "message", nil},
		{"reply comments", 1, 1, "message", nil},
		{"write comments with 0 product id", 0, 0, "message", errors.ErrCreateEntity},
	}

	for _, tc := range cases {
		payload := api.WriteCommentRequest{
			Productid: int32(tc.prodid),
			Replyid:   int32(tc.replyid),
			Message:   tc.message,
			Owner:     "user",
		}
		_, err := svc.WriteComment(ctx, &payload)
		assert.Equal(t, err, tc.err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
	}
}
