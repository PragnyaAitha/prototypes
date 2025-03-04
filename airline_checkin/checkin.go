// filepath: /workspaces/prototypes/airline_checkin/checkin.go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Database connection
	dsn := "root:Mysql@123@tcp(127.0.0.1:3306)/mysql"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Verify connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to the database!")

	var wg sync.WaitGroup

	for i := 0; i < 180; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			fmt.Printf("Goroutine %d\n", index)
			// Example database operation
			// _, err := db.Exec("INSERT INTO seats (user_id) VALUES (0)")
			// if err != nil {
			// 	log.Println(err)
			// }
			result, err := db.Exec("UPDATE seats SET user_id = ? WHERE user_id != 0 ORDER BY id ASC LIMIT 1", i)
			// fmt.Printf("Goroutine %d\n", index)
			if err != nil {
				log.Println(err)
			}
			rowsAffected, err := result.RowsAffected()
			if err != nil {
				log.Printf("Error fetching rows affected for goroutine %d: %v\n", index, err)
				return
			}
			if rowsAffected == 0 {
				log.Printf("No rows updated for goroutine %d\n", index)
			} else {
				log.Printf("Goroutine %d updated %d rows\n", index, rowsAffected)
			}
		}(i)
	}

	wg.Wait()
}
