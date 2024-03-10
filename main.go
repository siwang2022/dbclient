package main

import (
	"context"
	"log"
	"time"
)

func triggerSwitch(db *DbClient) {
	for range time.Tick(time.Second * 5) {
		log.Println("Time to switch db ")
		_ = db.SwitchDB()
	}
}

func doInsertMust(db *DbClient) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Failed to insert and recovered", r)
		}
		time.Sleep(1 * time.Second)
	}()

	tx := db.MustBegin()
	tx.MustExec("INSERT INTO person (first_name, last_name, email) VALUES ($1, $2, $3)", "a", "b", "c@jmoiron.net")
	time.Sleep(15 * time.Second)
	err := tx.Commit()
	if err != nil {
		log.Println("Commit failed:", err)
		return
	}

	log.Println("Committed a tx")
}

func doInsert(db *DbClient) {
	defer time.Sleep(1 * time.Second)

	ctx := context.Background() // we could use withDeadline to cancel the call
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		log.Println("BeginTxx failed:", err)
		return
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO person (first_name, last_name, email) VALUES ($1, $2, $3)", "a", "b", "c@jmoiron.net")
	if err != nil {
		log.Println("ExecContext failed:", err)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Println("Commit failed:", err)
		return
	}

	log.Println("Committed a tx")
}

func doQueryRows(db *DbClient) {
	type Person struct {
		FirstName string `db:"first_name"`
		LastName  string `db:"last_name"`
		Email     string
	}

	ctx := context.Background() // we could use withDeadline to cancel the call
	query := "select * from person"
	rows, err := db.QueryxContext(ctx, query)
	if err != nil {
		log.Println("QueryxContext failed:", err)
		return
	}

	for rows.Next() {
		person := Person{}
		err := rows.StructScan(&person)
		if err != nil {
			log.Println("rows.StructScan failed:", err)
			return
		}
	}

	log.Println("Query rows succeeded")
}

func main() {
	db := ConnectDB("postgresql://127.0.0.1:5432/swang")

	go triggerSwitch(db)

	for {
		// doInsertMust(db)
		doInsert(db)
		doQueryRows(db)
	}
}
