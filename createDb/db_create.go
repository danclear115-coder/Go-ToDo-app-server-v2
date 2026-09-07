package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

const databasePath = "../database/database.db"

func InitDB() {

	databaseDir := filepath.Dir(databasePath)

	if err := os.MkdirAll(databaseDir, 0755); err != nil {
		log.Fatal("create database directory:", err)
	}

	var err error

	DB, err = sql.Open("sqlite3", databasePath)
	if err != nil {
		log.Fatal("open database:", err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatal("ping database:", err)
	}

	if _, err := DB.Exec("PRAGMA foreign_keys = ON"); err != nil {
		log.Fatal("enable foreign keys:", err)
	}

	if err := createUsersTable(); err != nil {
		log.Fatal("create users table:", err)
	}

	if err := createTasksTable(); err != nil {
		log.Fatal("create tasks table:", err)
	}

	log.Println("Database is ready")
}
