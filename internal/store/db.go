package store

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx"
)

func PingDB(dataBase string) bool {

	db, err := sql.Open("pgx", dataBase)
	if err != nil {
		return false
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		return false
	}
	return true
}
