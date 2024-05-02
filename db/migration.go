package db

import (
	"database/sql"
	"strings"

	"github.com/soleimanim/gotype/logger"
)

type Migrator interface {
	Migrate() error
}

type MigrationHandler struct {
	dbHandler  *DBHandler
	migrations map[int]string
}

func (m *MigrationHandler) Migrate() error {
	var lastMigrationVersion int
	tableFound := true
	row, err := m.dbHandler.DB.Query("select max(version) from migrations")
	if err != nil && strings.HasPrefix(err.Error(), "no such table") {
		logger.Println("No migration table found. creating it...")
		lastMigrationVersion = -1
		tableFound = false
	} else if err != nil {
		logger.Println("Error migrating database", err)
		return err
	}

	if tableFound && row.Next() {
		err = row.Scan(&lastMigrationVersion)
		if err == sql.ErrNoRows {
			logger.Println("No migratoin record exists, executing all migrations")
			lastMigrationVersion = -1
		} else if err != nil {
			logger.Println("Error migrating database", err)
			return err
		}
	} else {
		lastMigrationVersion = -1
	}
	row.Close()

	for version, query := range m.migrations {
		if version > lastMigrationVersion {
			logger.Println("execuring migration version:", version)
			_, err := m.dbHandler.DB.Exec(query)
			if err != nil {
				logger.Println("Error executing migration", version, err)
				return err
			}
			_, err = m.dbHandler.DB.Exec("Insert into migrations (version) values (?)", version)
			if err != nil {
				logger.Println("Error inserting migration version to database", version, err)
			}
			lastMigrationVersion = version
			logger.Println("Migration version", version, "executed")
		}
	}

	logger.Println("Migration success")
	return nil
}

func NewMigrator(db *DBHandler) Migrator {
	mh := MigrationHandler{
		dbHandler: db,
		migrations: map[int]string{
			0: "CREATE TABLE IF NOT EXISTS migrations (" +
				"id INTEGER PRIMARY KEY AUTOINCREMENT," +
				"version INTEGER NOT NULL," +
				"migrated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP" +
				");",
			1: "CREATE TABLE IF NOT EXISTS " + TABLE_TYPING_TESTS + " (" +
				"id INTEGER PRIMARY KEY AUTOINCREMENT," +
				"test_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP," +
				"speed FLOAT NOT NULL," +
				"accuracy FLOAT NOT NULL," +
				"words_count INTEGER NOT NULL," +
				"mistakes_count INTEGER NOT NULL" +
				");",
		},
	}

	return &mh
}
