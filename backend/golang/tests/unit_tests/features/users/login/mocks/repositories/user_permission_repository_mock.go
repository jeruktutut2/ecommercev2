package mockrepositories

import (
	"backend-golang/features/users/login/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/mock"
)

type UserPermissionRepositoryMock struct {
	Mock mock.Mock
}

func (repository *UserPermissionRepositoryMock) FindByUserId(pool *pgxpool.Pool, ctx context.Context, userId int32) (userPermissions []models.UserPermission, err error) {
	arguments := repository.Mock.Called(pool, ctx, userId)
	return arguments.Get(0).([]models.UserPermission), arguments.Error(1)
}
