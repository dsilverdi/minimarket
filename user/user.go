package user

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type User struct {
	ID        string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Auth struct {
	Token string
	Email string
}

type JwtClaims struct {
	jwt.StandardClaims
	Email string
}

type UserRepository interface {
	Save(ctx context.Context, user User) error
	Read(ctx context.Context, email string) (*User, error)
}
