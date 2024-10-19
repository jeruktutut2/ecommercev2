package models

import "github.com/jackc/pgx/v5/pgtype"

type UserPermission struct {
	Id           pgtype.Int4
	UserId       pgtype.Int4
	PermissionId pgtype.Int4
}
