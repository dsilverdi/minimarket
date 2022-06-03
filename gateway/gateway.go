package gateway

import (
	"context"
	"minimarket/pkg/rest"
	productapi "minimarket/product/api"
	userapi "minimarket/user/api"
	"time"
)

type User struct {
	Email    string
	Password string
}

type Auth struct {
	Token string
}

type Product struct {
	ID                 int              `json:"id"`
	ProductName        string           `json:"product_name"`
	ProductCategory    string           `json:"product_category"`
	ProductPhoto       string           `json:"product_photo"`
	ProductDescription string           `json:"product_description"`
	ProductComment     []ProductComment `json:"product_comment"`
}

type ProductComment struct {
	ID        int       `json:"id"`
	ProductID int       `json:"product_id"`
	Message   string    `json:"message"`
	Owner     string    `json:"owner"`
	Replies   []Reply   `json:"replies"`
	CreatedAt time.Time `json:"created_at"`
}

type Reply struct {
	ID        int       `json:"id"`
	ProductID int       `json:"product_id"`
	Message   string    `json:"message"`
	Owner     string    `json:"owner"`
	CreatedAt time.Time `json:"created_at"`
}

type UserClientInterface interface {
	RegisterUser(ctx context.Context, email string, password string) (*userapi.UserIdentity, error)
	AuthorizeUser(ctx context.Context, email string, password string) (*userapi.AuthResponse, error)
}

type ProductClientInterface interface {
	GetProducts(ctx context.Context, name string, category string, query rest.QueryParameter) (*productapi.ProductResponse, error)
	WriteComment(ctx context.Context, productid int, replyid int, owner string, message string) (*productapi.CommentData, error)
}
