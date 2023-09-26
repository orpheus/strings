package main

import (
	"fmt"
	"github.com/orpheus/strings/infrastructure/postgres"
	"github.com/orpheus/strings/infrastructure/server"
	"github.com/orpheus/strings/util"
	"log"
)

func main() {
	log.Println("Starting server...")

	dbUser := util.GetEnv("DB_USER", "postgres")
	dbPass := util.GetEnv("DB_PASS", "")
	dbName := util.GetEnv("DB_NAME", "strings")
	dbHost := util.GetEnv("DB_HOST", "localhost")
	dbPort := util.GetEnv("DB_PORT", "5432")

	log.Println("Connecting database...")

	jdbcUrl := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	conn := postgres.NewPgxPool(jdbcUrl)
	defer conn.Close()

	log.Println("Running migrations...")

	postgres.Migrate(conn)

	log.Println("Creating gin server router...")

	s := server.NewGin()
	server.Construct(s, conn)

	log.Println("Running server...")
	err := s.Run()
	if err != nil {
		log.Fatalln("Error starting server")
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
