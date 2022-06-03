package gateway_test

import (
	"context"
	"fmt"
	"minimarket/gateway"
	"minimarket/gateway/mocks"
	"minimarket/pkg/errors"
	"minimarket/pkg/rest"
	"minimarket/product/api"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	validEmail    = "user@example.com"
	invalidEmail  = "userexample.com"
	validPass     = "password"
	invalidPass   = "wrong"
	existEmail    = "exist@example.com"
	existPassword = "existpass"
	contentType   = "application/json"
)

var (
	user      = gateway.User{Email: validEmail, Password: validPass}
	existUser = gateway.User{Email: existEmail, Password: existPassword}
)

func NewService() gateway.Service {
	commentlist := []*api.CommentData{
		{Id: 1, Productid: 1, Message: "comment", Owner: existEmail, Replyid: 0, Createdat: time.Now().Format(time.RFC3339)},
		{Id: 2, Productid: 1, Message: "comment", Owner: existEmail, Replyid: 1, Createdat: time.Now().Format(time.RFC3339)},
		{Id: 3, Productid: 1, Message: "comment", Owner: existEmail, Replyid: 1, Createdat: time.Now().Format(time.RFC3339)},
		{Id: 4, Productid: 2, Message: "comment", Owner: existEmail, Replyid: 0, Createdat: time.Now().Format(time.RFC3339)},
		{Id: 5, Productid: 2, Message: "comment", Owner: existEmail, Replyid: 0, Createdat: time.Now().Format(time.RFC3339)},
	}

	productlist := []*api.ProductData{
		{Id: 1, Name: "bajumock", Category: "pakaian", Picture: "google.com", Description: "this is deskripsi basic", Createdat: time.Now().Format(time.RFC3339), Updatedat: time.Now().Format(time.RFC3339)},
		{Id: 2, Name: "bajumock", Category: "pakaian", Picture: "google.com", Description: "this is deskripsi basic", Createdat: time.Now().Format(time.RFC3339), Updatedat: time.Now().Format(time.RFC3339)},
	}

	productmock := mocks.NewProductClient(productlist, commentlist)
	usermock := mocks.NewUserClient(existUser)
	return gateway.New(usermock, productmock)
}

func TestRegister(t *testing.T) {
	svc := NewService()
	ctx := context.Background()
	cases := []struct {
		desc     string
		email    string
		password string
		err      error
	}{
		{"register new user", validEmail, validPass, nil},
		{"register with existing user", existEmail, existPassword, errors.ErrAlreadyExists},
	}

	for _, tc := range cases {
		err := svc.Register(ctx, tc.email, tc.password)
		assert.Equal(t, err, tc.err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
	}
}

func TestAuthorize(t *testing.T) {
	svc := NewService()
	ctx := context.Background()
	cases := []struct {
		desc     string
		email    string
		password string
		err      error
	}{
		{"authorize new user", validEmail, validPass, errors.ErrNotFound},
		{"authorize existing user", existEmail, existPassword, nil},
		{"authorize wrong password", existEmail, "wrong", errors.ErrWrongPassword},
	}

	for _, tc := range cases {
		_, err := svc.Authorize(ctx, tc.email, tc.password)
		assert.Equal(t, err, tc.err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
	}
}

func TestGetProducts(t *testing.T) {
	svc := NewService()
	ctx := context.Background()
	cases := []struct {
		desc     string
		name     string
		category string
		err      error
	}{
		{"get product", "", "", nil},
		{"get product", "bajumock", "", nil},
		{"get product", "", "pakaian", nil},
	}

	for _, tc := range cases {
		_, err := svc.GetProducts(ctx, rest.QueryParameter{Limit: 10, Page: 1}, tc.name, tc.category)
		assert.Equal(t, err, tc.err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
	}
}

func TestWriteComment(t *testing.T) {
	svc := NewService()
	ctx := context.Background()
	cases := []struct {
		desc      string
		productid int
		replyid   int
		message   string
		owner     string
		err       error
	}{
		{"write new comment", 1, 1, "new comment", existEmail, nil},
		{"write new comment with 0 productid", 0, 1, "new comment", existEmail, errors.ErrCreateEntity},
	}

	for _, tc := range cases {
		err := svc.WriteComment(ctx, tc.productid, tc.replyid, tc.message, tc.owner)
		assert.Equal(t, err, tc.err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
	}
}
