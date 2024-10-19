package mockrepositories

import (
	"backend-golang/features/users/login/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	Mock mock.Mock
}

func (repository *UserRepositoryMock) FindByEmail(pool *pgxpool.Pool, ctx context.Context, email string) (user models.User, err error) {
	arguments := repository.Mock.Called(pool, ctx, email)
	return arguments.Get(0).(models.User), arguments.Error(1)
}
