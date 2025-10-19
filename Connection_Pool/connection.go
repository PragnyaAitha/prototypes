package main

import (
	"database/sql"
	"log"
	"sync"
	"time"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// var arr [10]*sql.DB
	// min := 1
	// max := 5

	// Define the data source name (DSN)
	dsn := "root:pass1234@tcp(localhost:3306)/mysql"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(10);
	db.SetMaxIdleConns(5)

	if err := db.Ping(); err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}

	start := time.Now()
	var wg sync.WaitGroup
	numWorkers := 10000
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			_, err := db.Exec("SELECT SLEEP(0.01)")
			if err != nil {
				log.Println("Query error:", err)
			}
		}(i)
	}
	wg.Wait()
	
	fmt.Println("Job done", time.Since(start))

}
