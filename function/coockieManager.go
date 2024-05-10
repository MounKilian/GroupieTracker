package groupieTracker

import (
	"errors"
	"log"
	"net/http"
	"strconv"
)

// Get a coockie from a specifique coockie name
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

func GetCoockieCode(w http.ResponseWriter, r *http.Request, name string) string {
	cookie, err := r.Cookie(name)
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			return ""
		default:
			log.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
	}
	code := cookie.Value
	return code
}

// Set user id inside a coockie
func SetCookie(w http.ResponseWriter, user User) {
	cookie := http.Cookie{
		Name:     "userId",
		Value:    strconv.Itoa(user.id),
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)
}

func SetCookieCode(w http.ResponseWriter, user User, code string) {
	cookie := http.Cookie{
		Name:     "code",
		Value:    code,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: false,
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
