package groupieTracker

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	_ "github.com/mattn/go-sqlite3"
)

type Error struct {
	errorMessage string
}

func Formulaire(w http.ResponseWriter, r *http.Request) User {
	db, err := sql.Open("sqlite3", "BDD.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
	}

	var user User
	for key, _ := range r.Form {
		if key == "pseudo-sign" {
			pseudo := r.FormValue("pseudo-sign")
			email := r.FormValue("email-sign")
			password := r.FormValue("password-sign")
			verifyPassword := r.FormValue("verify-password-sign")
			passwordEncrypt := Encrypt(password)
			valueCreate := [3]string{pseudo, email, passwordEncrypt}

			if verifyPassword != password && VerifyPassword(password) {
				log.Println("ERROR : Incorrect password")
			} else {
				createUser(db, valueCreate)
				valueConnect := [2]string{pseudo, passwordEncrypt}
				user = connectUser(db, valueConnect)
				if user.pseudo == "" {
					log.Println("ERROR : Wrong connection information")
				} else {
					http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
					return user
				}
			}
		} else if key == "connect-log" {
			password := r.FormValue("password-log")
			passwordEncrypt := Encrypt(password)
			connect := r.FormValue("connect-log")
			valueConnect := [2]string{connect, passwordEncrypt}

			user = connectUser(db, valueConnect)
			if user.pseudo == "" {
				log.Println("ERROR : Wrong pseudo")
			} else {
				http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
				return user
			}
		}
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
	return user
}

func VerifyPassword(s string) bool {
	var hasNumber, hasUpperCase, hasLowercase, hasSpecial bool
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsUpper(c):
			hasUpperCase = true
		case unicode.IsLower(c):
			hasLowercase = true
		case c == '#' || c == '|':
			return false
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}
	return hasNumber && hasUpperCase && hasLowercase && hasSpecial
}

func BlindForm(w http.ResponseWriter, r *http.Request, Blindtest Blindtest) {
	db, err := sql.Open("sqlite3", "BDD.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
	}
	musicName := r.FormValue("music-name")
	if MatchTitle(Blindtest.currentBtest.name, musicName) {
	}
	if Blindtest.currentPlay == Blindtest.nbSong {
		Blindtest.currentPlay = 0
		http.Redirect(w, r, "/win", http.StatusFound)
	}
	http.Redirect(w, r, "/blindtest", http.StatusFound)
}

func MatchTitle(title, input string) bool {
	if strings.Contains(title, " - ") {
		if index := strings.Index(title, " - "); index != -1 {
			title = title[:index]
		} else {
			title = ""
		}
	} else if strings.Contains(title, " (") {
		if index := strings.Index(title, " ("); index != -1 {
			title = title[:index]
		} else {
			title = ""
		}
	}
	if strings.ToLower(title) == strings.ToLower(input) {
		return true
	}
	return false
}

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
