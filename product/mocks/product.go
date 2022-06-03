package mocks

import (
	"context"
	"minimarket/pkg/errors"
	"minimarket/product"
)

type ProductRepository struct {
	query []product.ProductQuery
}

func NewProductRepository(products []product.ProductQuery) product.ProductRepository {
	return &ProductRepository{
		query: products,
	}
}

func (pr *ProductRepository) Read(ctx context.Context, name string, category string, limit int, page int) ([]product.ProductQuery, error) {

	var result []product.ProductQuery
	for _, q := range pr.query {
		if name != "" && q.ProductName == name {
			result = append(result, q)
		}

		if category != "" && q.ProductCategory == category {
			result = append(result, q)
		}
	}

	if len(result) < limit {
		limit = len(result)
	}

	var resultWithLimit []product.ProductQuery
	for i := 0; i < limit; i++ {
		resultWithLimit = append(resultWithLimit, result[i])
	}

	return resultWithLimit, nil
}

func (pr *ProductRepository) WriteComment(ctx context.Context, productID int, replyID int, message string, owner string) error {
	if productID == 0 {
		return errors.ErrCreateEntity
	}

	return nil
}
