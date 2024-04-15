package main

import (
	"database/sql"
	groupieTracker "groupieTracker/function"
	"log"
)

func main() {
	db, err := sql.Open("sqlite3", "BDD.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	groupieTracker.CreateRoomTest(db)
}
