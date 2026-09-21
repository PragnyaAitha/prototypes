package main

import (
	"context"
	"database/sql"
	"fmt"
)

// User is the row stored on exactly one shard. ID is the shard key.
type User struct {
	ID    int64
	Name  string
	Email string
}

// ShardRouter owns one connection pool per shard and applies the shard topology
// for every operation.
type ShardRouter struct {
	plan ShardPlan
	dbs  []*sql.DB
}

func NewShardRouter(ctx context.Context, rootDSN string) (*ShardRouter, error) {
	plan := NewShardPlan()
	router := &ShardRouter{plan: plan, dbs: make([]*sql.DB, len(plan.Partitions))}

	for index, shard := range plan.Partitions {
		db, err := sql.Open("mysql", shardDSN(rootDSN, shard.Database))
		if err != nil {
			router.Close()
			return nil, fmt.Errorf("open shard %s: %w", shard.Database, err)
		}
		if err := db.PingContext(ctx); err != nil {
			db.Close()
			router.Close()
			return nil, fmt.Errorf("ping shard %s: %w", shard.Database, err)
		}
		router.dbs[index] = db
	}

	return router, nil
}

func (router *ShardRouter) Close() error {
	var firstErr error
	for _, db := range router.dbs {
		if db == nil {
			continue
		}
		if err := db.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func (router *ShardRouter) ShardForKey(key int64) Shard {
	index := int(key % int64(len(router.plan.Partitions)))
	if index < 0 {
		index += len(router.plan.Partitions)
	}
	return router.plan.Partitions[index]
}

func (router *ShardRouter) shardDB(key int64) *sql.DB {
	index := int(key % int64(len(router.dbs)))
	if index < 0 {
		index += len(router.dbs)
	}
	return router.dbs[index]
}

func (router *ShardRouter) CreateUser(ctx context.Context, user User) error {
	if user.ID < 1 {
		return fmt.Errorf("user ID must be positive")
	}
	_, err := router.shardDB(user.ID).ExecContext(ctx,
		"INSERT INTO users (id, name, email) VALUES (?, ?, ?)",
		user.ID, user.Name, user.Email)
	if err != nil {
		return fmt.Errorf("insert user %d into %s: %w", user.ID, router.ShardForKey(user.ID).Database, err)
	}
	return nil
}

func (router *ShardRouter) GetUser(ctx context.Context, userID int64) (User, error) {
	var user User
	err := router.shardDB(userID).QueryRowContext(ctx,
		"SELECT id, name, email FROM users WHERE id = ?", userID).
		Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return User{}, fmt.Errorf("get user %d from %s: %w", userID, router.ShardForKey(userID).Database, err)
	}
	return user, nil
}

// ListUsers demonstrates a scatter-gather read for data without a shard key.
func (router *ShardRouter) ListUsers(ctx context.Context) ([]User, error) {
	users := make([]User, 0)
	for index, db := range router.dbs {
		rows, err := db.QueryContext(ctx, "SELECT id, name, email FROM users ORDER BY id")
		if err != nil {
			return nil, fmt.Errorf("list users from %s: %w", router.plan.Partitions[index].Database, err)
		}
		for rows.Next() {
			var user User
			if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
				rows.Close()
				return nil, err
			}
			users = append(users, user)
		}
		if err := rows.Err(); err != nil {
			rows.Close()
			return nil, err
		}
		rows.Close()
	}
	return users, nil
}
