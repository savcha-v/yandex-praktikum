package store

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/caarlos0/env"
)

var FileName string

type UnitURL struct {
	Full  string `json:"full"`
	Short string `json:"short"`
}

var urls = make(map[int]UnitURL)

type Config struct {
	FileStor string `env:"FILE_STORAGE_PATH" envDefault:""`
}

func InitStorage() {

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}
	FileName = cfg.FileStor
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

func GetID(urlToShort string) string {

	nextID := getNextID()
	until := UnitURL{
		Full:  urlToShort,
		Short: "fff",
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
	return strconv.Itoa(nextID)

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

func getNextID() int {
	return len(urls)
}
