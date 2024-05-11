package main

import (
	"database/sql"
	"fmt"
	groupieTracker "groupieTracker/function"
	"log"
)

func main() {
	db, err := sql.Open("sqlite3", "BDD.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	fmt.Println(groupieTracker.GetUsersScoreInRoom(db, 38))
}
