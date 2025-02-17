package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Define the data source name (DSN)
	dsn := "root:Mysql@1234@tcp(localhost:3306)/mysql"

	// // Open a connection to the database
	// db, err := sql.Open("mysql", dsn)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()

	// // Ping the database to verify the connection
	// err = db.Ping()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("Successfully connected to the database!")

	var arr [10]*sql.DB
	min := 0
	max := 5

	db, err := sql.Open("mysql", dsn)
	arr = append(arr, db)
	if len(arr) > 0 {
		conn := arr[0]
		conn.Exec("SELECT SLEEP(5)")
	} else {
		db, err := sql.Open("mysql", dsn)
		db.Exec("SELECT SLEEP(5)")
	}

}
