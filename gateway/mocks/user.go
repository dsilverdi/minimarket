package mocks

import (
	"context"
	"crypto/sha256"
	"fmt"
	"minimarket/gateway"
	"minimarket/pkg/errors"
	"minimarket/user/api"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var JWT_SECRET_KEY = []byte("minimarket-signature-key")

type UserClient struct {
	Existingusers gateway.User
}

type JwtClaims struct {
	jwt.StandardClaims
	Email string
}

func NewUserClient(existUser gateway.User) gateway.UserClientInterface {
	return UserClient{
		Existingusers: existUser,
	}
}

func (cl UserClient) RegisterUser(ctx context.Context, email string, password string) (*api.UserIdentity, error) {
	if email == cl.Existingusers.Email {
		return nil, errors.ErrAlreadyExists
	}
	return nil, nil
}

func (cl UserClient) AuthorizeUser(ctx context.Context, email string, password string) (*api.AuthResponse, error) {
	if email != cl.Existingusers.Email {
		return nil, errors.ErrNotFound
	}

	userpwd := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
	existpwd := fmt.Sprintf("%x", sha256.Sum256([]byte(cl.Existingusers.Password)))
	if userpwd != existpwd {
		return nil, errors.ErrWrongPassword
	}

	claims := &JwtClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(3 * 24 * time.Hour).Unix(),
		},
		Email: email,
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	tokenString, err := token.SignedString(JWT_SECRET_KEY)
	if err != nil {
		return nil, err
	}

	return &api.AuthResponse{
		Token: tokenString,
		Email: email,
	}, nil
}
