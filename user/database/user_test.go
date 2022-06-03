package database_test

import (
	"context"
	"database/sql"
	"fmt"
	"minimarket/pkg/errors"
	"minimarket/pkg/uuid"
	"minimarket/user"
	"minimarket/user/database"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var idProvider = uuid.New()

func TestSave(t *testing.T) {
	email := "user-save@mail.com"
	uid, err := idProvider.ID()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	cases := []struct {
		desc string
		user user.User
		err  error
	}{
		{
			desc: "new user",
			user: user.User{
				ID:        uid,
				Email:     email,
				Password:  "pass",
				CreatedAt: time.Now(),
			},
			err: nil,
		},
		{
			desc: "duplicate user",
			user: user.User{
				ID:        uid,
				Email:     email,
				Password:  "pass",
				CreatedAt: time.Now(),
			},
			err: errors.ErrCreateEntity,
		},
		{
			desc: "empty email",
			user: user.User{
				ID:        uid,
				Email:     "",
				Password:  "pass",
				CreatedAt: time.Now(),
			},
			err: errors.ErrCreateEntity,
		},
		{
			desc: "empty password",
			user: user.User{
				ID:        uid,
				Email:     email,
				Password:  "",
				CreatedAt: time.Now(),
			},
			err: errors.ErrCreateEntity,
		},
	}

	dbMiddleware := database.NewDatabase(db)
	repo := database.NewUsersRepository(dbMiddleware)

	for _, tc := range cases {
		err := repo.Save(context.Background(), tc.user)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestRead(t *testing.T) {
	dbMiddleware := database.NewDatabase(db)
	repo := database.NewUsersRepository(dbMiddleware)

	email := "user-retrieval@example.com"

	uid, err := idProvider.ID()
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	user := user.User{
		ID:        uid,
		Email:     email,
		Password:  "pass",
		CreatedAt: time.Now(),
	}

	err = repo.Save(context.Background(), user)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	cases := map[string]struct {
		email string
		err   error
	}{
		"existing user":     {email, nil},
		"non-existing user": {"unknown@example.com", sql.ErrNoRows},
	}

	for desc, tc := range cases {
		_, err := repo.Read(context.Background(), tc.email)
		assert.True(t, errors.Contains(err, tc.err), fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
	}
}
