package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/airline-voucher-seat-assignment/utils"
	_ "github.com/mattn/go-sqlite3"
)

func InitDB(config *utils.BaseConfig) (*sql.DB, error) {
	dbPath := config.DBPath
	if dbPath == "" {
		dbPath = "./vouchers.db"
	}

	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %v", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite database: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	if err = runSchema(db); err != nil {
		return nil, fmt.Errorf("failed to run schema: %v", err)
	}

	if config.Environment == "development" {
		fmt.Printf("SQLite database initialized successfully at: %s\n", dbPath)
	}

	return db, nil
}

func runSchema(db *sql.DB) error {
	schemaPath := "./database/schema.sql"

	schema, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		schema = []byte(`
		CREATE TABLE IF NOT EXISTS vouchers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			crew_name TEXT NOT NULL,
			crew_id TEXT NOT NULL,
			flight_number TEXT NOT NULL,
			flight_date TEXT NOT NULL,
			aircraft_type TEXT NOT NULL,
			seat1 TEXT NOT NULL,
			seat2 TEXT NOT NULL,
			seat3 TEXT NOT NULL,
			created_at TEXT NOT NULL
		);`)
	}

	if _, err := db.Exec(string(schema)); err != nil {
		return fmt.Errorf("failed to execute schema: %v", err)
	}

	return nil
}
