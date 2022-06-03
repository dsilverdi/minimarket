package mocks

import (
	"context"
	"minimarket/pkg/errors"
	"minimarket/user"
)

type UserRepository struct {
	existuser user.User
}

func NewUsersRepository(user user.User) user.UserRepository {
	return &UserRepository{
		existuser: user,
	}
}

func (u *UserRepository) Save(ctx context.Context, user user.User) error {
	if user.Email == "" || user.Password == "" {
		return errors.ErrCreateEntity
	}

	if u.existuser.Email == user.Email {
		return errors.New("Duplicate entry")
	}
	return nil
}

func (u *UserRepository) Read(ctx context.Context, email string) (*user.User, error) {
	if email != u.existuser.Email {
		return nil, errors.ErrNotFound
	}

	return &u.existuser, nil
}
