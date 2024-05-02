package db

import (
	"database/sql"
	"os"
	"path/filepath"
	"runtime"

	_ "github.com/mattn/go-sqlite3"

	"github.com/jmoiron/sqlx"
	"github.com/soleimanim/gotype/logger"
)

type DBHandler struct {
	DB *sqlx.DB
}

func NewDBHandler() *DBHandler {
	return &DBHandler{}
}

// Open a new connection to the sqlite database
//
// Returns:
//   - error: will return error struct or nil in case of success
func (db *DBHandler) Open() error {
	dbPath, err := databasePath()
	if err != nil {
		logger.Println("Error getting database path", err)
		return err
	}

	logger.Println("Opening database at path", dbPath)
	db.DB, err = sqlx.Open("sqlite3", dbPath)
	if err != nil {
		logger.Println("Error initializing the database", err)
		return err
	}

	err = db.DB.Ping()
	if err != nil {
		return err
	}

	// Run database migrations
	mh := NewMigrator(db)
	err = mh.Migrate()
	if err != nil {
		return err
	}

	return nil
}

func (db *DBHandler) Close() {
	db.DB.Close()
}

func (db *DBHandler) Exec(query string, args ...any) (sql.Result, error) {
	return db.DB.Exec(query, args...)
}

// path of sqlite3 database
//
// Returns:
//   - string: path to the database file
//   - error: error of reading the path or nil
func databasePath() (string, error) {
	dbPath := ""
	dbName := "data"
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	switch os := runtime.GOOS; os {
	case "darwin":
		dbPath = filepath.Join(homeDir, "Library", "Application Support", "GoType")
	case "linux":
		dbPath = filepath.Join(homeDir, ".gotype")
	case "windows":
		dbPath = filepath.Join(homeDir, "AppData", "Local", "GoType")
	}

	_, err = os.Stat(dbPath)
	if os.IsNotExist(err) {
		os.MkdirAll(dbPath, 0700)
	} else if err != nil {
		return "", err
	}
	dbPath = filepath.Join(dbPath, dbName)

	return dbPath, nil
}
