package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver
)

var DB *sql.DB

func InitDB(config *Config) {
	connStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		config.DBUser,
		config.DBPassword,
		config.DBName,
		config.DBHost,
		config.DBPort,
		config.DBSSLMode,
	)

	var err error
	DB, err = sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal("❌ Failed to open DB:", err)
	}

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = DB.PingContext(ctx)
	if err != nil {
		log.Fatal("❌ Failed to ping DB:", err)
	}

	log.Println("✅ Connected to PostgreSQL!")

	// Setup table
	setupSQL := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        name TEXT NOT NULL,
        email TEXT NOT NULL UNIQUE
    );

    DELETE FROM users;
    INSERT INTO users (id, name, email) VALUES (1, 'Alice', 'alice@example.com');
    INSERT INTO users (id, name, email) VALUES (2, 'Bob', 'bob@example.com');
    `

	_, err = DB.Exec(setupSQL)
	if err != nil {
		log.Fatal("❌ Failed to setup table:", err)
	}
}
