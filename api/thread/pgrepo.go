package thread

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/orpheus/strings/api"
	"github.com/orpheus/strings/core"
	"github.com/orpheus/strings/infrastructure/logging"
	"log"
)

type Repository struct {
	DB     api.PgxConn
	Logger logging.Logger
}

func (r *Repository) FindAll() ([]core.Thread, error) {
	threadRows, err := r.DB.Query(context.Background(), "select * from thread")
	if err != nil {
		return nil, err
	}
	defer threadRows.Close()

	var threads []core.Thread
	for threadRows.Next() {
		r := new(core.Thread)
		err := threadRows.Scan(&r.Id, &r.Name, &r.Description, &r.DateCreated, &r.DateModified)
		if err != nil {
			log.Fatal(err)
		}
		threads = append(threads, *r)
	}

	if err := threadRows.Err(); err != nil {
		return nil, err
	}

	if len(threads) == 0 {
		return []core.Thread{}, nil
	}

	return threads, nil
}

func (r *Repository) CreateOne(thread core.Thread) (core.Thread, error) {
	sql := "insert into thread (name, description) " +
		"VALUES ($1, $2) " +
		"RETURNING id, name, description"
	t := new(core.Thread)
	err := r.DB.QueryRow(context.Background(), sql, thread.Name, thread.Description).
		Scan(&t.Id, &t.Name, &t.Description)

	fmt.Println(r, err)
	return *t, err
}

func (r *Repository) DeleteById(id uuid.UUID) error {
	sql := "delete from thread where id = $1"
	_, err := r.DB.Exec(context.Background(), sql, id)
	return err
}
