package cookie

import (
	"net/http"
	"yandex-praktikum/internal/config"
	"yandex-praktikum/internal/encryption"

	uuid "github.com/google/uuid"
)

func SetUserID(cfg config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// func SetUserID(next http.Handler) http.Handler {
			// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// получим куки для идентификации пользователя
			_, err := r.Cookie("userID")
			if err != nil {
				// если не нашли - добавим новую
				userID := uuid.New().String()
				dst, err := encryption.Encrypt(userID, cfg)
				if err != nil {
					dst = string(userID)
				}
				newCookie := &http.Cookie{
					Name:  "userID",
					Value: dst,
				}
				http.SetCookie(w, newCookie)
				r.AddCookie(newCookie)
			}
			next.ServeHTTP(w, r)
		})
	}
}

func GetUserID(r *http.Request, cfg config.Config) string {
	userID := ""
	if cookieUserID, err := r.Cookie("userID"); err == nil {
		userID, err = encryption.Decrypt(cookieUserID.Value, cfg)
		if err != nil {
			userID = cookieUserID.Value
		}
	}
	return userID
}
