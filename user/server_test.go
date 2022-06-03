package user_test

import (
	"context"
	"crypto/sha256"
	"fmt"
	"minimarket/pkg/errors"
	"minimarket/pkg/uuid"
	"minimarket/user"
	"minimarket/user/api"
	"minimarket/user/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	existEmail = "exist@mail.com"
	existPass  = "exist@pass.com"
	validEmail = "user@mail.com"
	validPass  = "user@mail.com"
)

func NewService() api.AuthServiceServer {
	usermocks := mocks.NewUsersRepository(user.User{
		Email:    existEmail,
		Password: fmt.Sprintf("%x", sha256.Sum256([]byte(existPass))),
	})
	uuid := uuid.New()
	return user.NewServer(usermocks, uuid)
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
		{"register existing user", existEmail, existPass, errors.ErrAlreadyExists},
		{"register empty email", "", existPass, errors.ErrCreateEntity},
		{"register empty password", validEmail, "", errors.ErrCreateEntity},
	}

	for _, tc := range cases {
		payload := api.UserRequest{
			Email:    tc.email,
			Password: tc.password,
		}
		_, err := svc.Register(ctx, &payload)
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
		{"authorize non-existing user", validEmail, validPass, errors.ErrNotFound},
		{"authorize existing user", existEmail, existPass, nil},
		{"authorize wrong password", existEmail, "wrongpass", errors.ErrWrongPassword},
	}

	for _, tc := range cases {
		payload := api.UserRequest{
			Email:    tc.email,
			Password: tc.password,
		}
		_, err := svc.Authorize(ctx, &payload)
		assert.Equal(t, err, tc.err, fmt.Sprintf("%s: unexpected error %s", tc.desc, err))
	}
}
