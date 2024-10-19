package models

import "github.com/jackc/pgx/v5/pgtype"

type User struct {
	Id        pgtype.Int4
	Username  pgtype.Text
	Email     pgtype.Text
	Password  pgtype.Text
	CreatedAt pgtype.Int8
}
