package store

import (
	"context"
	"log"
	"strconv"
	"sync"
	config "yandex-praktikum/internal/config"
)

type unitURL struct {
	Full   string `json:"full"`
	Short  string `json:"short"`
	UserID string `json:"userID"`
}

type UserShorts struct {
	Short string `json:"short_url"`
	Full  string `json:"original_url"`
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

func GetShortURL(urlToShort string, host string, cfg config.Config, userID string) string {

	mu := &sync.Mutex{}
	mu.Lock()
	defer mu.Unlock()

	until := unitURL{
		Full:   urlToShort,
		Short:  "http://" + host + "/" + cfg.BaseURL + "/" + "?id=",
		UserID: userID,
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
	return until.Short
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
		for _, unitURL := range urls {
			if unitURL.UserID != userID {
				continue
			}
			unitRes := UserShorts{
				Short: unitURL.Short,
				Full:  unitURL.Full,
			}
			result = append(result, unitRes)
		}
	}
	return result
}
