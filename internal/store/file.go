package store

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
)

func fileInit(fileStor string) error {
	file, err := os.OpenFile(fileStor, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		if scanner.Err() != nil {
			return scanner.Err()
		}
		data := scanner.Bytes()

		json.Unmarshal([]byte(data), &urls)
	}
	return nil
}

func fileWrite(fileStor string, until unitURL, nextID int) error {
	file, err := os.OpenFile(fileStor, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
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
	return nil
}
