package store

import (
	"strconv"
)

var urls = make(map[int]string)

func GetID(urlToShort string) string {

	nextID := getNextID()
	urls[nextID] = urlToShort
	return strconv.Itoa(nextID)

}

func GetURL(id string) (url string, strErr string) {

	idStr, err := strconv.Atoi(id)
	if err != nil {
		return "", "'id' not found"
	}

	url, exists := urls[idStr]
	if !exists {
		return "", "'id' not found"
	}

	return url, ""
}

func getNextID() int {
	return len(urls)
}
