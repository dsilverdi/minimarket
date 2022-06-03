package database

import (
	"context"
	"minimarket/pkg/errors"
	"minimarket/user"
	"time"
)

type UserRepository struct {
	db Database
}
type UserDB struct {
	ID        string    `db:"id"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
}

func NewUsersRepository(db Database) user.UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) Save(ctx context.Context, user user.User) error {
	query := `INSERT INTO USERS (id, email, password, created_at)
	VALUES (:id, :email, :password, :created_at);`

	userDB := &UserDB{
		ID:        user.ID,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
	}

	_, err := u.db.NamedExecContext(ctx, query, userDB)
	if err != nil {
		return errors.Wrap(errors.ErrCreateEntity, err)
	}

	return nil
}

func (u *UserRepository) Read(ctx context.Context, email string) (*user.User, error) {
	var userDB UserDB

	query := `SELECT id, email, password, created_at FROM USERS WHERE email = ?`

	err := u.db.QueryRowxContext(ctx, query, email).StructScan(&userDB)
	if err != nil {
		return nil, err
	}

	User := &user.User{
		ID:        userDB.ID,
		Email:     userDB.Email,
		Password:  userDB.Password,
		CreatedAt: userDB.CreatedAt,
	}

	return User, nil
}
