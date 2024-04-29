package groupieTracker

import (
	"errors"
	"log"
	"net/http"
	"strconv"
)

func GetCoockie(w http.ResponseWriter, r *http.Request, name string) int {
	cookie, err := r.Cookie(name)
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			http.Error(w, "cookie not found", http.StatusBadRequest)
		default:
			log.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
	}
	userId, _ := strconv.Atoi(cookie.Value)
	return userId
}

func SetCookie(w http.ResponseWriter, user User) {
	cookie := http.Cookie{
		Name:     "userId",
		Value:    strconv.Itoa(user.id),
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)
}

func DeleteCookies(w http.ResponseWriter, r *http.Request) {
	cookies := r.Cookies()
	for _, cookie := range cookies {
		cookie.MaxAge = -1
		cookie.Secure = false
		http.SetCookie(w, cookie)
	}
}
