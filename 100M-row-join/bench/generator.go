package main

import (
	"bytes"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"sync"
)

// env: CONNECTION

func main() {

	connection := os.Getenv("CONNECTION")
	threads := 32

	db, err := sql.Open("mysql", connection)
	if err != nil {
		fmt.Printf("Could not connect to MySQL: %s.", err)
	}

	db.SetMaxOpenConns(threads)
	w := new(sync.WaitGroup)
	w.Add(1)

	fmt.Printf("Running go rountines\n")
	abuf := new(bytes.Buffer)
	bbuf := new(bytes.Buffer)

	for a := 0; a <= 10000; a++ {
		// assume that b_id is mostly null
		if a < 1 {
			abuf.WriteString("(REPLACE(UUID(),'-',''), REPLACE(UUID(),'-','')),")
		} else {
			abuf.WriteString("(REPLACE(UUID(),'-',''), NULL),")
		}
		bbuf.WriteString("(REPLACE(UUID(),'-','')),")
	}
	tblA := fmt.Sprintf("INSERT INTO a (id, b_id) VALUES %s (REPLACE(UUID(),'-',''), REPLACE(UUID(),'-',''));", abuf.String())
	tblB := fmt.Sprintf("INSERT INTO b (id) VALUES %s (REPLACE(UUID(),'-',''));", bbuf.String())

	for i := 0; i < threads; i++ {
		go insertOnLoop(db, tblA)
		go insertOnLoop(db, tblB)
	}

	fmt.Printf("Waiting for go routines\n")
	w.Wait()
}

func insertOnLoop(db *sql.DB, sql string) {

	for {
		_, err := db.Exec(sql)
		if err != nil {
			fmt.Printf("error: %s\n", err)
		} else {
			fmt.Printf(".")
		}
	}
}
