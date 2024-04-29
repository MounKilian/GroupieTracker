package groupieTracker

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	id       int
	pseudo   string
	email    string
	password string
}

func selectFromTable(db *sql.DB, table string) *sql.Rows {
	query := "SELECT * FROM " + table
	result, _ := db.Query(query)
	return result
}

func checkUser(db *sql.DB, value [3]string) int {
	var nbAccount int

	query := "SELECT COUNT(*) FROM USER WHERE pseudo = ? OR email = ?"
	err := db.QueryRow(query, value[0], value[1]).Scan(&nbAccount)
	if err != nil {
		log.Fatal(err)
	}
	return nbAccount
}

func createUser(db *sql.DB, value [3]string) {
	if checkUser(db, value) == 0 {
		insertQuery := "INSERT INTO USER (id, pseudo, email, password) VALUES (?, ?, ?, ?)"
		_, err := db.Exec(insertQuery, nil, value[0], value[1], value[2])
		if err != nil {
			log.Fatal(err)
		}
	}
}

func connectUser(db *sql.DB, value [2]string) User {
	var u User
	db.QueryRow("SELECT * FROM `USER` WHERE pseudo = ? OR email = ? AND password = ?", value[0], value[0], value[1]).Scan(&u.id, &u.pseudo, &u.email, &u.password)
	return u
}

func displayUserRows(rows *sql.Rows) {
	for rows.Next() {
		var u User
		err := rows.Scan(&u.id, &u.pseudo, &u.email, &u.password)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(u)
	}
}

func resetUserTable(db *sql.DB) {
	_, err := db.Exec("DROP TABLE IF EXISTS `USER`")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("CREATE TABLE USER (id INTEGER PRIMARY KEY, pseudo TEXT NOT NULL, email TEXT NOT NULL, password TEXT NOT NULL)")
	if err != nil {
		log.Fatal(err)
	}
}

func GetUserById(db *sql.DB, id int) User {
	var u User
	db.QueryRow("SELECT * FROM `USER` WHERE id = ?", id).Scan(&u.id, &u.pseudo, &u.email, &u.password)
	return u
}
