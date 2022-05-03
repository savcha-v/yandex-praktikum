package encryption

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
)

var key = "10c57de0"

func Decrypt(msg string) (string, error) {

	// выделяем подпись
	dst := msg[:len(msg)-36]
	// выделяем id
	id := strings.Replace(msg, dst, "", -1)
	// декодируем в hex
	data, err := hex.DecodeString(dst)
	if err != nil {
		panic(err)
	}
	// хеш
	h := hmac.New(sha256.New, []byte(key))
	// вычисляем подпись
	h.Write([]byte(id))
	sign := h.Sum(nil)
	// Проверить подпись
	if hmac.Equal(data, sign) {
		return id, nil
	} else {
		return "", errors.New("incorrect userID")
	}
}

func Encrypt(src string) (string, error) {

	data := []byte(src)
	// вычисляем хеш
	h := hmac.New(sha256.New, []byte(key))
	h.Write(data)
	dst := hex.EncodeToString(h.Sum(nil))
	msg := dst + src
	return msg, nil
}
