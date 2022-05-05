package store

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
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

func InitStorage(fileStor string) {

	if fileStor != "" {
		file, err := os.OpenFile(fileStor, os.O_RDONLY|os.O_CREATE, 0777)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {

			if scanner.Err() != nil {
				log.Fatal(scanner.Err())
			}
			data := scanner.Bytes()

			if err != nil {
				log.Fatal(err)
			}
			json.Unmarshal([]byte(data), &urls)
		}
	}
}

func GetShortURL(urlToShort string, host string, cfg config.Config, userID string) string {

	mu := &sync.Mutex{}
	mu.Lock()
	defer mu.Unlock()

	nextID := len(urls)

	shortURL := "http://" + host + "/" + cfg.BaseURL + "/" + "?id=" + strconv.Itoa(nextID)

	until := unitURL{
		Full:   urlToShort,
		Short:  shortURL,
		UserID: userID,
	}
	urls[nextID] = until

	//записать в файл
	if cfg.FileStor != "" {

		file, err := os.OpenFile(cfg.FileStor, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		var record = make(map[int]unitURL)

		record[nextID] = until
		data, err := json.Marshal(record)
		if err != nil {
			log.Fatal(err)
		}

		writer := bufio.NewWriter(file)
		if _, err := writer.Write(data); err != nil {
			log.Fatal(err)
		}
		if err := writer.WriteByte('\n'); err != nil {
			log.Fatal(err)
		}
		writer.Flush()
	}
	return shortURL
}

func GetURL(idStr string) (url string, strErr string) {

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return "", "'id' not found"
	}

	until, exists := urls[id]
	if !exists {
		return "", "'id' not found"
	}

	return until.Full, ""
}

func GetUserShorts(userID string) []UserShorts {

	var result []UserShorts
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
	return result
}
