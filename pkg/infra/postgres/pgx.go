package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgtype"
	pgtypeuuid "github.com/jackc/pgtype/ext/gofrs-uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"os"
)

// NewPgxPool creates a new Postgres connection pool that satisfies
// the PgxConn interface inside the repository package for use
// with our repos.
func NewPgxPool(jdbcUrl string) *pgxpool.Pool {
	dbConfig, err := pgxpool.ParseConfig(jdbcUrl)
	if err != nil {
		log.Fatalln("Could not parse pgx connection string")
	}

	conn, err := pgxpool.ConnectConfig(context.Background(), dbConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	dbConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		conn.ConnInfo().RegisterDataType(pgtype.DataType{
			Value: &pgtypeuuid.UUID{},
			Name:  "uuid",
			OID:   pgtype.UUIDOID,
		})
		return nil
	}

	return conn
}
