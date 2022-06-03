package database_test

import (
	"context"
	"fmt"
	"minimarket/pkg/errors"
	"minimarket/product/database"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadQuery(t *testing.T) {
	dbMiddleware := database.NewDatabase(db)
	repo := database.NewProductRepository(dbMiddleware)

	err := save(context.Background(), dbMiddleware, "sample-baju", "pakaian", "description", "url")
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	err = save(context.Background(), dbMiddleware, "sample-baju", "pakaian", "description", "url")
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	//comment
	for i := 0; i < 5; i++ {
		err = repo.WriteComment(context.Background(), 2, 0, "message-test", "user")
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
	}

	//reply
	for i := 0; i < 5; i++ {
		err = repo.WriteComment(context.Background(), 2, 2, "message-test", "user")
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
	}

	cases := map[string]struct {
		prodname     string
		prodcategory string
		limit        int
		page         int
		err          error
	}{
		"read data":                 {"", "", 15, 1, nil},
		"read with search name":     {"sample-baju", "", 15, 1, nil},
		"read with search category": {"", "pakaian", 15, 1, nil},
		"read with pagination":      {"", "", 15, 2, nil},
	}

	for desc, tc := range cases {
		_, err := repo.Read(context.Background(), tc.prodname, tc.prodcategory, tc.limit, tc.page)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
	}
}

func TestWriteComment(t *testing.T) {
	dbMiddleware := database.NewDatabase(db)
	repo := database.NewProductRepository(dbMiddleware)

	cases := map[string]struct {
		prodid  int
		repid   int
		message string
		owner   string
		err     error
	}{
		"write comment":                   {1, 0, "message", "user", nil},
		"reply comment":                   {1, 1, "message", "user", nil},
		"write comment with 0 product id": {0, 0, "message", "user", errors.ErrCreateEntity},
	}

	for desc, tc := range cases {
		err := repo.WriteComment(context.Background(), tc.prodid, tc.repid, tc.message, tc.owner)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
	}
}

func save(ctx context.Context, db database.Database, name, category, description, picture string) error {
	query := `INSERT INTO PRODUCT (name, description, category, picture, created_at, updated_at)
	VALUES (:name, :description, :category, :picture, :created_at, :updated_at);`

	product := map[string]interface{}{
		"name":        name,
		"description": description,
		"category":    category,
		"picture":     picture,
		"created_at":  time.Now(),
		"updated_at":  time.Now(),
	}

	_, err := db.NamedExecContext(ctx, query, product)
	if err != nil {
		return err
	}

	return nil
}
