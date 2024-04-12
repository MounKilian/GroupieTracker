package groupieTracker

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

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
			passwordEncrypt := Encrypt(password)
			valueCreate := [3]string{pseudo, email, passwordEncrypt}

			createUser(db, valueCreate)
			valueConnect := [2]string{pseudo, passwordEncrypt}
			user = connectUser(db, valueConnect)
			if user.pseudo == "" {
				log.Println("ERROR : Wrong connection information")
			} else {
				http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
				return user
			}
		} else if key == "connect-log" {
			password := r.FormValue("password-log")
			passwordEncrypt := Encrypt(password)
			connect := r.FormValue("connect-log")
			valueConnect := [2]string{connect, passwordEncrypt}

			user = connectUser(db, valueConnect)
			if user.pseudo == "" {
				log.Println("ERROR : Wrong connection information")
			} else {
				http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
				return user
			}
		}
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
	return user
}
