package repositories

import (
	"backend-golang/features/users/login/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserPermissionRepository interface {
	FindByUserId(pool *pgxpool.Pool, ctx context.Context, userId int32) (userPermissions []models.UserPermission, err error)
}

type UserPermissionRepositoryImplementation struct {
}

func NewUserPermissinoRepository() UserPermissionRepository {
	return &UserPermissionRepositoryImplementation{}
}

func (repository *UserPermissionRepositoryImplementation) FindByUserId(pool *pgxpool.Pool, ctx context.Context, userId int32) (userPermissions []models.UserPermission, err error) {
	rows, err := pool.Query(ctx, `SELECT id, user_id, permission_id FROM user_permissions WHERE user_id = $1;`, userId)
	if err != nil {
		return
	}
	defer func() {
		rows.Close()
		if rows.Err() != nil {
			userPermissions = []models.UserPermission{}
			err = rows.Err()
		}
	}()

	for rows.Next() {
		var userPermission models.UserPermission
		err = rows.Scan(&userPermission.Id, &userPermission.UserId, &userPermission.PermissionId)
		if err != nil {
			userPermissions = []models.UserPermission{}
			return
		}
		userPermissions = append(userPermissions, userPermission)
	}
	return
}
