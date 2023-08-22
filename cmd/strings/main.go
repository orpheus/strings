package main

import (
	"fmt"
	"github.com/orpheus/strings/pkg/infrastructure/postgres"
	"github.com/orpheus/strings/pkg/infrastructure/server"
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
	dbUser := getEnv("DB_USER", "postgres")
	dbPass := getEnv("DB_PASS", "")
	dbName := getEnv("DB_NAME", "strings-v2")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")

	jdbcUrl := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	conn := postgres.NewPgxPool(jdbcUrl)
	defer conn.Close()
	postgres.Migrate(conn)

	s := server.NewGin()
	server.Construct(s, conn)

	err := s.Run()
	if err != nil {
		panic("Error starting server")
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
