package initialize

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateTableUser(pool *pgxpool.Pool, ctx context.Context) {
	query := `CREATE TABLE users (
  		id SERIAL PRIMARY KEY,
  		username varchar(50) NOT NULL UNIQUE,
  		email varchar(100) NOT NULL UNIQUE,
  		password varchar(100) NOT NULL,
  		created_at bigint NOT NULL
	);`
	_, err := pool.Exec(ctx, query)
	if err != nil {
		log.Fatalln("error when creating table users:", err.Error())
	}
	log.Println("create table user succedded")
}

func CreateDataUser(pool *pgxpool.Pool, ctx context.Context) {
	query := `INSERT INTO users (username,email,password,created_at) VALUES ('username','email@email.com','$2a$10$MvEM5qcQFk39jC/3fYzJzOIy7M/xQiGv/PAkkoarCMgsx/rO0UaPG',1695095017);`
	_, err := pool.Exec(ctx, query)
	if err != nil {
		log.Fatalln("error when creating data users:", err.Error())
	}
	log.Println("create data user succedded")
}

type User struct {
	Id        pgtype.Int4
	Username  pgtype.Text
	Email     pgtype.Text
	Password  pgtype.Text
	CreatedAt pgtype.Int8
}

func GetDataUserByEmail(pool *pgxpool.Pool, ctx context.Context, email string) (user User) {
	query := `SELECT id, username, email, password, created_at FROM users WHERE email = $1;`
	err := pool.QueryRow(ctx, query, email).Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		log.Fatalln("error when getting data users:", err.Error())
	}
	log.Println("get data user succedded")
	return
}

func DropTableUser(pool *pgxpool.Pool, ctx context.Context) {
	query := `DROP TABLE IF EXISTS users;`
	_, err := pool.Exec(ctx, query)
	if err != nil {
		log.Fatalln("error when dropping table users:", err.Error())
	}
	log.Println("drop table user succedded")
}
