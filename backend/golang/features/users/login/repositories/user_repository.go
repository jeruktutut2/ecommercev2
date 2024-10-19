package repositories

import (
	"backend-golang/features/users/login/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	FindByEmail(pool *pgxpool.Pool, ctx context.Context, email string) (user models.User, err error)
}

type UserRepositoryImplementation struct {
}

func NewUserRepository() UserRepository {
	return &UserRepositoryImplementation{}
}

func (repository *UserRepositoryImplementation) FindByEmail(pool *pgxpool.Pool, ctx context.Context, email string) (user models.User, err error) {
	err = pool.QueryRow(ctx, `SELECT id, username, email, password, created_at FROM users WHERE email = $1;`, email).Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	return
}
