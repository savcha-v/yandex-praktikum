package store

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"strconv"
	config "yandex-praktikum/internal/config"
)

var FileName string

type UnitURL struct {
	Full  string `json:"full"`
	Short string `json:"short"`
}

var urls = make(map[int]UnitURL)

func InitStorage() {

	FileName = config.Cfg.FileStor

	if FileName != "" {

		file, err := os.OpenFile(FileName, os.O_RDONLY|os.O_CREATE, 0777)
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

func GetShortURL(urlToShort string, host string) string {

	nextID := len(urls)

	shortURL := "http://" + host + "/" + config.Cfg.BaseURL + "/" + "?id=" + strconv.Itoa(nextID)

	until := UnitURL{
		Full:  urlToShort,
		Short: shortURL,
	}
	urls[nextID] = until
	//записать в файл
	if FileName != "" {

		file, err := os.OpenFile(FileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		var record = make(map[int]UnitURL)

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
