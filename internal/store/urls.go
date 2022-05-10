package store

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"sync"
	config "yandex-praktikum/internal/config"
)

type unitURL struct {
	Full       string `json:"full"`
	Short      string `json:"short"`
	UserID     string `json:"userID"`
	UUID       string
	httpStatus int
}

type UserShorts struct {
	Short string `json:"short_url"`
	Full  string `json:"original_url"`
}

type responseURL struct {
	UUID  string `json:"correlation_id"`
	Short string `json:"short_url"`
}

type RequestURL struct {
	Full   string `json:"original_url,omitempty"`
	Short  string
	UserID string
	UUID   string `json:"correlation_id,omitempty"`
}

var urls = make(map[int]unitURL)

func InitStorage(cfg config.Config) {

	if cfg.DataBase != "" {
		// context.Background(),
		if err := dbInit(cfg.DataBase); err != nil {
			log.Fatal(err)
		}
	} else if cfg.FileStor != "" {
		if err := fileInit(cfg.FileStor); err != nil {
			log.Fatal(err)
		}
	}
}

func GetShortURL(urlToShort string, host string, cfg config.Config, userID string) (string, int) {

	mu := &sync.Mutex{}
	mu.Lock()
	defer mu.Unlock()

	until := unitURL{
		Full:       urlToShort,
		Short:      "http://" + host + "/" + cfg.BaseURL + "/" + "?id=",
		UserID:     userID,
		httpStatus: http.StatusCreated,
	}

	if cfg.DataBase != "" {
		// записать в базу данных
		err := dbWrite(context.Background(), cfg.DataBase, &until)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// записать в память
		nextID := len(urls)
		until.Short += strconv.Itoa(nextID)
		urls[nextID] = until

		if cfg.FileStor != "" {
			//записать в файл
			if err := fileWrite(cfg.FileStor, until, nextID); err != nil {
				log.Fatal(err)
			}
		}

	}
	return until.Short, until.httpStatus
}

func ShortURLs(urls []RequestURL, host string, cfg config.Config, userID string) []responseURL {

	mu := &sync.Mutex{}
	mu.Lock()
	defer mu.Unlock()

	shortBase := "http://" + host + "/" + cfg.BaseURL + "/" + "?id="

	db, err := sql.Open("pgx", cfg.DataBase)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// объявляем транзакцию
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	// если возникает ошибка, откатываем изменения
	defer tx.Rollback()

	ctx := context.Background()

	// готовим инструкцию
	stmt, err := tx.PrepareContext(ctx, `INSERT INTO urls ("ID", "Full", "Short", "UserID") VALUES ($1, $2, $3, $4)`)
	if err != nil {
		log.Fatal(err)
	}
	// не забываем закрыть инструкцию, когда она больше не нужна
	defer stmt.Close()

	initialID, err := dbCountUrls(ctx, db)
	if err != nil {
		log.Fatal(err)
	}

	var result []responseURL

	for counter, url := range urls {
		nextID := initialID + counter
		short := shortBase + strconv.Itoa(nextID)
		if _, err = stmt.ExecContext(ctx, nextID, url.Full, short, userID); err != nil {
			log.Fatal(err)
		}

		res := responseURL{
			UUID:  url.UUID,
			Short: short,
		}
		result = append(result, res)

	}
	// сохраняем изменения
	tx.Commit()

	return result
}

func GetURL(ctx context.Context, idStr string, cfg config.Config) (url string, strErr string) {

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return "", "'id' not found"
	}

	var fullURL string
	if cfg.DataBase != "" {
		full, err := dbReadURL(ctx, cfg.DataBase, id)
		if err != nil {
			return "", err.Error()
		}
		fullURL = full
	} else {
		until, exists := urls[id]
		if !exists {
			return "", "'id' not found"
		}
		fullURL = until.Full
	}

	return fullURL, ""
}

func GetUserShorts(ctx context.Context, cfg config.Config, userID string) []UserShorts {

	var result []UserShorts
	if cfg.DataBase != "" {
		result = dbReadUserShorts(cfg.DataBase, userID)
	} else {
		for _, UnitURL := range urls {
			if UnitURL.UserID != userID {
				continue
			}
			unitRes := UserShorts{
				Short: UnitURL.Short,
				Full:  UnitURL.Full,
			}
			result = append(result, unitRes)
		}
	}
	return result
}
