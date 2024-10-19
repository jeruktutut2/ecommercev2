package initialize

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateTablePermission(pool *pgxpool.Pool, ctx context.Context) {
	query := `CREATE TABLE permissions (
  		id SERIAL PRIMARY KEY,
  		permission varchar(50) NOT NULL UNIQUE
	);`
	_, err := pool.Exec(ctx, query)
	if err != nil {
		log.Fatalln("error when creating table permissions:", err.Error())
	}
	log.Println("create table permissions succedded")
}

func CreateDataPermission(pool *pgxpool.Pool, ctx context.Context) {
	query := `INSERT INTO permissions (id, permission) VALUES (1, 'ADMINISTRATOR'), (2, 'CREATE_PERMISSION'), (3, 'READ_PERMISSION'), (4, 'UPDATE_PERMISSION'), (5, 'DELETE_PERMISSION');`
	_, err := pool.Exec(ctx, query)
	if err != nil {
		log.Fatalln("error when creating data permissions:", err.Error())
	}
	log.Println("create data permissions succedded")
}

func DropTablePermission(pool *pgxpool.Pool, ctx context.Context) {
	query := `DROP TABLE IF EXISTS permissions;`
	_, err := pool.Exec(ctx, query)
	if err != nil {
		log.Fatalln("error when dropping table permissions:", err.Error())
	}
	log.Println("drop table permissions succedded")
}
