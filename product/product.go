package product

import (
	"context"
	"time"
)

type ProductQuery struct {
	ProductID          int32      `db:"product_id"`
	ProductName        string     `db:"product_name"`
	ProductDescription string     `db:"product_description"`
	ProductCategory    string     `db:"product_category"`
	ProductPicture     string     `db:"product_picture"`
	CommentID          *int32     `db:"comment_id"`
	CommentMessage     *string    `db:"message"`
	CommentAuthor      *string    `db:"owner"`
	ReplyID            *int32     `db:"reply_id"`
	CreatedAt          *time.Time `db:"created_at"`
}

type CommentDB struct {
	ID        int       `db:"id"`
	ProductID int       `db:"product_id"`
	Message   string    `db:"message"`
	ReplyID   int       `db:"reply_id"`
	Owner     string    `db:"owner"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type ProductRepository interface {
	Read(ctx context.Context, name string, category string, limit int, page int) ([]ProductQuery, error)
	WriteComment(ctx context.Context, productID int, ReplyID int, message string, owner string) error
}

// type CommentRepository interface {
// 	Read(ctx context.Context, name string, category string, limit int, page int) ([]CommentDB, error)
// }
