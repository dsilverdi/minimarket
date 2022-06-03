package user

import (
	"context"
	"crypto/sha256"
	"fmt"
	"minimarket"
	"minimarket/pkg/errors"
	"minimarket/user/api"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var JWT_SECRET_KEY = []byte("minimarket-signature-key")

type UserServer struct {
	User       UserRepository
	IDProvider minimarket.IDprovider
	api.UnimplementedAuthServiceServer
}

func NewServer(user UserRepository, idprov minimarket.IDprovider) api.AuthServiceServer {
	return &UserServer{
		User:       user,
		IDProvider: idprov,
	}
}

func (svc *UserServer) Register(ctx context.Context, req *api.UserRequest) (*api.UserIdentity, error) {
	// var mysqlErr *mysql.MySQLError
	if req.Password == "" {
		return nil, errors.ErrCreateEntity
	}

	NewUser := &User{
		Email:     req.Email,
		Password:  fmt.Sprintf("%x", sha256.Sum256([]byte(req.Password))),
		CreatedAt: time.Now(),
	}

	id, err := svc.IDProvider.ID()
	if err != nil {
		return nil, errors.Wrap(errors.ErrCreateUUID, err)
	}

	NewUser.ID = id

	// Perform DB Call Here
	err = svc.User.Save(ctx, *NewUser)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return nil, errors.ErrAlreadyExists
		}
		return nil, errors.ErrCreateEntity
	}

	return &api.UserIdentity{Email: req.Email}, nil
}

func (svc *UserServer) Authorize(ctx context.Context, req *api.UserRequest) (*api.AuthResponse, error) {
	CurrentUser, err := svc.User.Read(ctx, req.Email)
	if err != nil {
		return nil, errors.ErrNotFound
	}

	userpwd := fmt.Sprintf("%x", sha256.Sum256([]byte(req.Password)))
	if userpwd != CurrentUser.Password {
		return nil, errors.ErrWrongPassword
	}

	claims := &JwtClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(3 * 24 * time.Hour).Unix(),
		},
		Email: CurrentUser.Email,
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
		Email: CurrentUser.Email,
	}, nil
}
