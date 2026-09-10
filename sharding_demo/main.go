package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// ShardPlan models a database partitioned into 3 logical shards.
// Each shard has its own database and tables for users, profiles,
// and a user_profiles junction table.
type ShardPlan struct {
	Partitions []Shard
}

type Shard struct {
	Database string
	Tables   []string
}

func NewShardPlan() ShardPlan {
	return ShardPlan{
		Partitions: []Shard{
			{
				Database: "app_shard_01",
				Tables:   []string{"users", "profiles", "user_profiles"},
			},
			{
				Database: "app_shard_02",
				Tables:   []string{"users", "profiles", "user_profiles"},
			},
			{
				Database: "app_shard_03",
				Tables:   []string{"users", "profiles", "user_profiles"},
			},
		},
	}
}

func main() {
	plan := NewShardPlan()
	fmt.Println("Shard plan ready:")
	for _, shard := range plan.Partitions {
		fmt.Printf("- database=%s tables=%v\n", shard.Database, shard.Tables)
	}

	rootDSN := os.Getenv("MYSQL_DSN")
	if rootDSN == "" {
		fmt.Println("No MYSQL_DSN is set. Showing shard plan only. Set MYSQL_DSN to create the databases and tables.")
		return
	}

	for _, shard := range plan.Partitions {
		createDatabaseIfNeeded(rootDSN, shard.Database)
		createSchema(rootDSN, shard)
	}

	fmt.Println("Successfully created 3 sharded partitions:", len(plan.Partitions))
}

func createDatabaseIfNeeded(dsn, dbName string) {
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	_, err = conn.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
	if err != nil {
		log.Fatal(err)
	}
}

func createSchema(rootDSN string, shard Shard) {
	shardDSN := shardDSN(rootDSN, shard.Database)
	conn, err := sql.Open("mysql", shardDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	_, err = conn.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			name VARCHAR(100) NOT NULL,
			email VARCHAR(150) NOT NULL UNIQUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Exec(`
		CREATE TABLE IF NOT EXISTS profiles (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			user_id BIGINT NOT NULL,
			bio TEXT,
			city VARCHAR(80),
			country VARCHAR(80),
			FOREIGN KEY (user_id) REFERENCES users(id)
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Exec(`
		CREATE TABLE IF NOT EXISTS user_profiles (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			user_id BIGINT NOT NULL,
			profile_id BIGINT NOT NULL,
			role VARCHAR(30) DEFAULT 'member',
			UNIQUE KEY uq_user_profile (user_id, profile_id),
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (profile_id) REFERENCES profiles(id)
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Created tables in shard:", shard.Database)
}

func shardDSN(rootDSN, dbName string) string {
	if strings.Contains(rootDSN, "/mysql?") {
		return strings.Replace(rootDSN, "/mysql?", "/"+dbName+"?", 1)
	}
	if strings.Contains(rootDSN, "/mysql") {
		return strings.Replace(rootDSN, "/mysql", "/"+dbName, 1)
	}
	return rootDSN + "/" + dbName + "?parseTime=true"
}
