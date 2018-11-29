package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"sync"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)


func main() {

	db, err := sql.Open("mysql", "root@tcp(localhost:4000)/test")
	if err != nil {
		zap.S().Fatalf("Could not connect to MySQL: %s.", err)
	}

	db.SetMaxOpenConns(32)
	w := new(sync.WaitGroup)
	w.Add(1)

	fmt.Printf("Running go rountines\n")
	abuf := new(bytes.Buffer)
	bbuf := new(bytes.Buffer)

	// todo: for TiDB do we need to pad out the row to be comparable?

	for a := 0; a <= 1000; a++ {
		// assume that b_id is mostly null
		if a < 3 {
			abuf.WriteString("(REPLACE(UUID(),'-',''), REPLACE(UUID(),'-','')),")
		} else {
			abuf.WriteString("(REPLACE(UUID(),'-',''), NULL),")
		}
		bbuf.WriteString("(REPLACE(UUID(),'-','')),")
	}
	tblA := fmt.Sprintf("INSERT INTO a (id, b_id) VALUES %s (REPLACE(UUID(),'-',''), REPLACE(UUID(),'-',''));", abuf.String())
	tblB := fmt.Sprintf("INSERT INTO b (id) VALUES %s (REPLACE(UUID(),'-',''));", bbuf.String())

	for i := 0; i < 32; i++ {
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

