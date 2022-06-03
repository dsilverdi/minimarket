package api

import (
	"encoding/json"
	"fmt"
	"io"
	"minimarket/gateway"
	"minimarket/gateway/mocks"
	"minimarket/product/api"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
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
	tokenmock     = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTQyNjg2MTEsIkVtYWlsIjoiZHNpbHZlcmRpQGdtYWlsLmNvbSJ9.t1Xgin102svDPQoWiQji8YasIlm3vTBZzmtqsndGWC0"
)

var (
	user      = UserRequestBody{Email: validEmail, Password: validPass}
	existUser = UserRequestBody{Email: existEmail, Password: existPassword}
)

type testRequest struct {
	method string
	url    string
	body   io.Reader
	token  string
}

func (tr testRequest) make(next echo.HandlerFunc) (*httptest.ResponseRecorder, error) {
	e := echo.New()
	req := httptest.NewRequest(tr.method, tr.url, tr.body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, tr.token)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	return rec, next(c)
}

func NewService() gateway.Service {
	commentlist := []*api.CommentData{
		{Id: 1, Productid: 1, Message: "comment basic", Owner: existEmail, Replyid: 0, Createdat: time.Now().Format(time.RFC3339)},
		{Id: 2, Productid: 1, Message: "comment basic", Owner: existEmail, Replyid: 0, Createdat: time.Now().Format(time.RFC3339)},
	}

	productlist := []*api.ProductData{
		{Id: 1, Name: "bajumock", Category: "pakaian", Picture: "google.com", Description: "this is deskripsi basic", Createdat: time.Now().Format(time.RFC3339), Updatedat: time.Now().Format(time.RFC3339)},
	}

	userdatamock := gateway.User{
		Email: existEmail, Password: existPassword,
	}

	productmock := mocks.NewProductClient(productlist, commentlist)
	usermock := mocks.NewUserClient(userdatamock)
	return gateway.New(usermock, productmock)
}

func toJSON(data interface{}) string {
	jsonData, _ := json.Marshal(data)
	return string(jsonData)
}

func TestRegisterEndpoint(t *testing.T) {

	svc := NewService()
	ep := NewEndpoint(svc)
	data := toJSON(user)
	withExistingUser := toJSON(existUser)
	withEmptyData := toJSON(gateway.User{})
	withEmptyEmail := toJSON(gateway.User{Password: validPass})
	withEmptyPassword := toJSON(gateway.User{Email: validEmail})

	cases := []struct {
		desc   string
		req    string
		status int
	}{
		{"register new user", data, http.StatusCreated},
		{"register with empty data", withEmptyData, http.StatusBadRequest},
		{"register with empty password", withEmptyPassword, http.StatusBadRequest},
		{"register with empty email", withEmptyEmail, http.StatusBadRequest},
		{"register with existing user", withExistingUser, http.StatusInternalServerError},
	}

	for _, tc := range cases {
		req := testRequest{
			method: http.MethodPost,
			url:    "/user/register",
			body:   strings.NewReader(tc.req),
		}

		rec, err := req.make(ep.RegisterEndpoint)
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		assert.Equal(t, tc.status, rec.Code, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, rec.Code))
	}
}

func TestAuthorizeEndpoint(t *testing.T) {
	svc := NewService()
	ep := NewEndpoint(svc)
	data := toJSON(existUser)
	withWrongPassword := toJSON(gateway.User{Email: existEmail, Password: invalidPass})
	withEmailNotFound := toJSON(gateway.User{Email: validEmail, Password: existPassword})
	withEmptyData := toJSON(gateway.User{})
	withEmptyEmail := toJSON(gateway.User{Password: validPass})
	withEmptyPassword := toJSON(gateway.User{Email: validEmail})

	cases := []struct {
		desc   string
		req    string
		status int
	}{
		{"authorize user", data, http.StatusOK},
		{"authorize with empty data", withEmptyData, http.StatusBadRequest},
		{"authorize with empty password", withEmptyPassword, http.StatusBadRequest},
		{"authorize with empty email", withEmptyEmail, http.StatusBadRequest},
		{"authorize with wrong password", withWrongPassword, http.StatusInternalServerError},
		{"authorize with Email Not Found", withEmailNotFound, http.StatusInternalServerError},
	}

	for _, tc := range cases {
		req := testRequest{
			method: http.MethodPost,
			url:    "/user/authorize",
			body:   strings.NewReader(tc.req),
		}

		rec, err := req.make(ep.AuthorizeEndpoint)
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		assert.Equal(t, tc.status, rec.Code, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, rec.Code))
	}
}

func TestGetProductEndpoint(t *testing.T) {
	svc := NewService()
	ep := NewEndpoint(svc)

	cases := []struct {
		desc   string
		url    string
		status int
	}{
		{"get products", "/product", http.StatusOK},
		{"get products", "/product?name=bajumock", http.StatusOK},
		{"get products", "/product?category=pakaian", http.StatusOK},
		{"get products", "/product?page=5&limit=1000", http.StatusOK},
	}

	for _, tc := range cases {
		req := testRequest{
			method: http.MethodGet,
			url:    tc.url,
			body:   nil,
		}

		rec, err := req.make(productCache(ep.GetProductEndpoint))
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		assert.Equal(t, tc.status, rec.Code, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, rec.Code))
	}
}

func TestWriteCommentEndpoint(t *testing.T) {
	svc := NewService()
	ep := NewEndpoint(svc)

	withNewData := toJSON(CommentRequestBody{ProductID: 1, ReplyID: 1, Message: "New Message"})
	withEmptyData := toJSON(CommentRequestBody{})
	cases := []struct {
		desc      string
		req       string
		authtoken string
		status    int
	}{
		{"write unauthorized", withNewData, "", http.StatusUnauthorized},
		{"write non valid token", withNewData, "Bearer ", http.StatusBadRequest},
		{"write comment", withNewData, tokenmock, http.StatusOK},
		{"write comment with empty data", withEmptyData, tokenmock, http.StatusBadRequest},
	}

	for _, tc := range cases {
		req := testRequest{
			method: http.MethodGet,
			url:    "/product/comment",
			body:   strings.NewReader(tc.req),
			token:  tc.authtoken,
		}

		rec, err := req.make(AuthMiddleware(ep.WriteCommentEndpoint))
		assert.Nil(t, err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
		assert.Equal(t, tc.status, rec.Code, fmt.Sprintf("%s: expected status code %d got %d", tc.desc, tc.status, rec.Code))
	}
}
