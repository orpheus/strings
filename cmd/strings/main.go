package main

import (
	"fmt"
	"github.com/orpheus/strings/pkg/infra/ginserver"
	"github.com/orpheus/strings/pkg/infra/sqldb"
	"github.com/orpheus/strings/pkg/infra/sqldb/migrations"
	"log"
	"os"
)

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func main() {
	sqlMigrationDir := getEnv("SQL_MIGRATION_DIR", "/Users/roark/code/github/orpheus/strings-go/pkg/infra/sqldb/migrations/sql")

	dbUser := getEnv("DB_USER", "postgres")
	dbPass := getEnv("DB_PASS", "")
	dbName := getEnv("DB_NAME", "stringsv2")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")

	store, err := sqldb.NewStore(sqldb.DriverPostgres, sqldb.ConnConfig{
		Host:   dbHost,
		Port:   dbPort,
		User:   dbUser,
		Pass:   dbPass,
		Dbname: dbName,
	})
	if err != nil {
		log.Fatalf("Failed creating pgx store: %s\n", err)
	}

	defer store.Close()

	err = migrations.Migrate(sqlMigrationDir, store.Db)
	if err != nil {
		log.Fatalf(fmt.Sprintf("failed to run migrations: %s\n", err.Error()))
	}

	s := ginserver.NewGin()
	ginserver.Construct(s, store)

	err = s.Run()
	if err != nil {
		log.Fatal("Error starting server")
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
