package migrations

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type migration struct {
	version  string
	filename string
	hash     string
}

// Migrate runs the migration scripts against the database instance.
// Migration script files are sql files that exist within the dir
// as pointed to by the variable sqlDir.
//
// Files must follow the flyway standard filename format:
//
//    VX.X.X__NAME
//
// Each file's hash is saved to a schema_history table and checked
// against each time migrations are ran. If a file that has already
// been migrated has changed, migration will fail.
//
// All migrations are committed in a transaction meaning if anything
// fails inside Migrate, the whole transaction is rolled back.
func Migrate(sqlDir string, db *sqlx.DB) error {
	return _migrate(sqlDir, db, nil)
}

// MigrateIgnore allows you to pass a list of filenames to ignore during
// migration. Useful for tests if you don't want the data population
// scripts to run.
func MigrateIgnore(sqlDir string, db *sqlx.DB, ignore []string) error {
	return _migrate(sqlDir, db, ignore)
}

func _migrate(sqlDir string, db *sqlx.DB, ignore []string) error {
	files, err := os.ReadDir(sqlDir)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	var txErr error

	defer func() {
		if txErr != nil {
			tx.Rollback()
			fmt.Printf("Rolled backed migration: %s\n", txErr)
		} else {
			tx.Commit()
		}
	}()

	// Setup migration table to track metadata
	migrationSql := `create table if not exists schema_history (
        version varchar primary key,
        filename varchar unique not null,
        hash varchar unique not null
    );`
	_, err = tx.Exec(migrationSql)
	if err != nil {
		return err
	}

	// Read in existing migrations
	rows, err := tx.Query(`select * from schema_history`)
	if err != nil {
		return err
	}
	defer rows.Close()
	migrations := make(map[string]migration)
	for rows.Next() {
		if rows.Err() != nil {
			return rows.Err()
		}
		var m migration
		err := rows.Scan(&m.version, &m.filename, &m.hash)
		if err != nil {
			return err
		}
		migrations[m.version] = m
	}

	fmt.Println("Beginning sql migrations...")
	for _, file := range files {
		filename := file.Name()

		// quick check to skip a migration file
		ignoreMigration := false
		for _, ignoreFile := range ignore {
			if filename == ignoreFile {
				ignoreMigration = true
				break
			}
		}
		if ignoreMigration {
			continue
		}

		regexStr := "^V\\d+.\\d+.\\d+__\\w+\\.(sql)$"
		validator := regexp.MustCompile(regexStr)

		validMigrationScript := validator.MatchString(filename)
		if !validMigrationScript {
			txErr = errors.New(fmt.Sprintf("Invalid migration script format: %s. Expecting %s\n", filename, regexStr))
			break
		}

		sqlBytes, err := os.ReadFile(filepath.Join(sqlDir, filename))
		if err != nil {
			txErr = err
			break
		}

		split := strings.Split(filename, "__")
		version := split[0]
		var existingMigration migration
		if _, ok := migrations[version]; ok {
			existingMigration = migrations[version]
		}

		// Generate sha hash of the file
		hasher := sha256.New()
		hasher.Write(sqlBytes)
		hash := hex.EncodeToString(hasher.Sum(nil))

		// If existing migration exist, check the hash hasn't changed
		if existingMigration != (migration{}) {
			if existingMigration.hash != hash {
				txErr = errors.New(fmt.Sprintf("%s file has changed, sha256 does not equal expected\n", filename))
				break
			}
			fmt.Printf("Skipping existing migration, %s...\n", filename)
			// if file was already migrated, then skip
			continue
		}

		sqlStr := string(sqlBytes)
		_, err = tx.Exec(sqlStr)
		if err != nil {
			txErr = errors.New(fmt.Sprintf("%s: %s\n", filename, err.Error()))
			break
		}

		// update migration table
		_, err = tx.Exec(`insert into schema_history (version, filename, hash) values ($1, $2, $3)`, version, filename, hash)
		if err != nil {
			txErr = errors.New(fmt.Sprintf("(%s) %s", filename, err.Error()))
			break
		}

		fmt.Printf("Ran migration file: %s\n", filename)
	}
	return txErr
}
