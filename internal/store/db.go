package store

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func PingDB(ctx context.Context, dataBase string) int {

	db, err := sql.Open("pgx", dataBase)
	if err != nil {
		return http.StatusInternalServerError
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

func dbInit(dataBase string) error {
	// ctx context.Context,
	db, err := sql.Open("pgx", dataBase)
	if err != nil {
		return err
	}
	defer db.Close()

	// ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	// defer cancel()

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS urls(
							"ID" INTEGER,
							"Full" TEXT PRIMARY KEY,
  							"Short" TEXT,
 							"UserID" TEXT
							 );`); err != nil {
		return err
	}
	return nil
}

func dbWrite(ctx context.Context, dataBase string, until *unitURL) error {

	db, err := sql.Open("pgx", dataBase)
	if err != nil {
		return err
	}
	defer db.Close()

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
	_, err = db.Exec(textInsert, id, until.Full, until.Short, until.UserID)

	if err != nil {
		pgErr := err.(*pgconn.PgError)

		if pgErr.Code == pgerrcode.UniqueViolation {
			short, err := dbReadShort(context.Background(), dataBase, until.Full)
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
	row := db.QueryRowContext(context.Background(), "SELECT COUNT(*) as count FROM urls")
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func dbReadURL(ctx context.Context, dataBase string, id int) (string, error) {

	db, err := sql.Open("pgx", dataBase)
	if err != nil {
		return "", err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var fullID sql.NullString

	textQuery := `SELECT "Full" FROM urls WHERE "ID" = $1`
	err = db.QueryRowContext(ctx, textQuery, id).Scan(&fullID)
	if err != nil {
		return "", err
	}

	if fullID.Valid {
		return fullID.String, nil
	}
	return "", errors.New("id not found")
}

func dbReadShort(ctx context.Context, dataBase string, full string) (string, error) {

	db, err := sql.Open("pgx", dataBase)
	if err != nil {
		return "", err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var fullID sql.NullString

	textQuery := `SELECT "Short" FROM urls WHERE "Full" = $1`
	err = db.QueryRowContext(ctx, textQuery, full).Scan(&fullID)
	if err != nil {
		return "", err
	}

	if fullID.Valid {
		return fullID.String, nil
	}
	return "", errors.New("full not found")
}

func dbReadUserShorts(dataBase string, userID string) []UserShorts {

	// ctx context.Context,
	db, err := sql.Open("pgx", dataBase)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	// defer cancel()

	textQuery := `SELECT "Full", "Short" FROM urls WHERE "UserID" = $1`
	rows, err := db.Query(textQuery, userID)
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
