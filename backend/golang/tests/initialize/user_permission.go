package initialize

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateTableUserPermission(pool *pgxpool.Pool, ctx context.Context) {
	query := `CREATE TABLE user_permissions (
  		id SERIAL PRIMARY KEY,
  		user_id int NOT NULL,
  		permission_id int NOT NULL,
    	CONSTRAINT user_permission_ibfk_1 FOREIGN KEY(user_id) REFERENCES users(id),
    	CONSTRAINT user_permission_ibfk_2 FOREIGN KEY(permission_id) REFERENCES permissions(id)
	);`
	_, err := pool.Exec(ctx, query)
	if err != nil {
		log.Fatalln("error when creating table user_permissions:", err.Error())
	}
	log.Println("create table user_permissions succedded")
}

func CreateDataUserPermission(pool *pgxpool.Pool, ctx context.Context) {
	query := `INSERT INTO user_permissions(user_id, permission_id) VALUES (1, 1), (1, 2), (1, 3), (1, 4), (1, 5);`
	_, err := pool.Exec(ctx, query)
	if err != nil {
		log.Fatalln("error when creating data user_permissions:", err.Error())
	}
	log.Println("create data user_permissions succedded")
}

func DropTableUserPermission(pool *pgxpool.Pool, ctx context.Context) {
	query := `DROP TABLE IF EXISTS user_permissions;`
	_, err := pool.Exec(ctx, query)
	if err != nil {
		log.Fatalln("error when dropping table user_permissions:", err.Error())
	}
	log.Println("drop table user_permissions succedded")
}
