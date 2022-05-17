package store

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"
	"yandex-praktikum/internal/config"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func PingDB(ctx context.Context, db *sql.DB) int {

	if err := db.PingContext(ctx); err != nil {
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

func dbInit(cfg *config.Config) error {

	db, err := sql.Open("pgx", cfg.DataBase)
	if err != nil {
		return err
	}

	textCreate := `CREATE TABLE IF NOT EXISTS urls(
		"ID" INTEGER,
		"Full" TEXT PRIMARY KEY,
		  "Short" TEXT,
		 "UserID" TEXT
		 );`
	if _, err := db.Exec(textCreate); err != nil {
		return err
	}

	cfg.ConnectDB = db

	return nil
}

func dbWrite(ctx context.Context, db *sql.DB, until *unitURL) error {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	id, err := dbCountUrls(ctx, db)
	if err != nil {
		return err
	}

	until.Short += strconv.Itoa(id)

	textInsert := `
	INSERT INTO urls ("ID", "Full", "Short", "UserID")
	VALUES ($1, $2, $3, $4)`
	_, err = db.ExecContext(ctx, textInsert, id, until.Full, until.Short, until.UserID)

	if err != nil {
		pgErr := err.(*pgconn.PgError)

		if pgErr.Code == pgerrcode.UniqueViolation {
			short, err := dbReadShort(ctx, db, until.Full)
			if err != nil {
				log.Fatal(err)
			}
			until.Short = short
			until.httpStatus = http.StatusConflict
			return nil
		}
		return err
	}

	return nil
}

func dbCountUrls(ctx context.Context, db *sql.DB) (int, error) {
	var id int

	// ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	// defer cancel()

	row := db.QueryRowContext(ctx, "SELECT COUNT(*) as count FROM urls")
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func dbReadURL(ctx context.Context, db *sql.DB, id int) (string, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var fullID sql.NullString

	textQuery := `SELECT "Full" FROM urls WHERE "ID" = $1`
	err := db.QueryRowContext(ctx, textQuery, id).Scan(&fullID)
	if err != nil {
		return "", err
	}

	if fullID.Valid {
		return fullID.String, nil
	}
	return "", errors.New("id not found")
}

func dbReadShort(ctx context.Context, db *sql.DB, full string) (string, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var fullID sql.NullString

	textQuery := `SELECT "Short" FROM urls WHERE "Full" = $1`
	err := db.QueryRowContext(ctx, textQuery, full).Scan(&fullID)
	if err != nil {
		return "", err
	}

	if fullID.Valid {
		return fullID.String, nil
	}
	return "", errors.New("full not found")
}

func dbReadUserShorts(ctx context.Context, db *sql.DB, userID string) []UserShorts {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	textQuery := `SELECT "Full", "Short" FROM urls WHERE "UserID" = $1`
	rows, err := db.QueryContext(ctx, textQuery, userID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var result []UserShorts

	for rows.Next() {
		var element UserShorts
		err = rows.Scan(&element.Full, &element.Short)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, element)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return result
}
