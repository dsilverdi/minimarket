package database

import (
	"context"
	"fmt"
	"minimarket/pkg/errors"
	"minimarket/product"
	"time"
)

type ProductRepository struct {
	db Database
}

func NewProductRepository(db Database) product.ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (pr *ProductRepository) Read(ctx context.Context, name string, category string, limit int, page int) ([]product.ProductQuery, error) {
	var productqueries []product.ProductQuery

	querymap := make(map[string]interface{})

	if name != "" {
		querymap["name"] = name
	}

	if category != "" {
		querymap["category"] = category
	}

	offset := (page - 1) * limit
	query_product := `SELECT id, name, description, category, picture FROM PRODUCT WHERE 1 `

	for key, val := range querymap {
		query_product += fmt.Sprintf(`AND %s = '%s' `, key, val)
	}

	query_product += fmt.Sprintf(`ORDER BY id LIMIT %d OFFSET %d`, limit, offset)

	query := fmt.Sprintf(`SELECT
		p.id as product_id, p.name as product_name, p.description as product_description, p.category as product_category, p.picture as product_picture,
		PRODUCT_COMMENT.id as comment_id, PRODUCT_COMMENT.message as message, PRODUCT_COMMENT.owner as owner, PRODUCT_COMMENT.reply_ID as reply_id,
		PRODUCT_COMMENT.created_at as created_at
		FROM (
			%s
			) as p
		LEFT JOIN 
		PRODUCT_COMMENT ON p.id = PRODUCT_COMMENT.product_id;`, query_product)

	rows, err := pr.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var productDB product.ProductQuery
		err = rows.StructScan(&productDB)
		if err != nil {
			return nil, err
		}

		productqueries = append(productqueries, productDB)
	}

	return productqueries, nil
}

func (pr *ProductRepository) WriteComment(ctx context.Context, productID int, replyID int, message string, owner string) error {
	query := `INSERT INTO PRODUCT_COMMENT (product_id, message, reply_id, owner, created_at, updated_at)
	VALUES (:product_id, :message, :reply_id, :owner, :created_at, :updated_at);`

	userDB := &product.CommentDB{
		ProductID: productID,
		Message:   message,
		ReplyID:   replyID,
		Owner:     owner,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := pr.db.NamedExecContext(ctx, query, userDB)
	if err != nil {
		return errors.Wrap(errors.ErrCreateEntity, err)
	}

	return nil
}
