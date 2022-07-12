package string

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/orpheus/strings/api"
	"github.com/orpheus/strings/core"
	"github.com/orpheus/strings/infrastructure/logging"
	"log"
)

type StringRepository struct {
	DB     api.PgxConn
	Logger logging.Logger
}

func (s *StringRepository) FindAll() ([]core.String, error) {
	rows, err := s.DB.Query(context.Background(), "select * from string")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var strings []core.String
	for rows.Next() {
		var r core.String
		err := rows.Scan(&r.Id, &r.Name, &r.Order, &r.Thread, &r.Description, &r.DateCreated, &r.DateModified)
		if err != nil {
			log.Fatal(err)
		}
		strings = append(strings, r)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	s.Logger.Logf("Fetched %d strings\n", len(strings))

	if len(strings) == 0 {
		return []core.String{}, nil
	}

	return strings, nil
}

func (s *StringRepository) FindAllByThread(threadId uuid.UUID) ([]core.String, error) {
	sql := "select * from string where thread = $1 order by \"order\" asc"
	rows, err := s.DB.Query(context.Background(), sql, threadId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var strings []core.String
	for rows.Next() {
		var r core.String
		err := rows.Scan(&r.Id, &r.Name, &r.Order, &r.Thread, &r.Description, &r.DateCreated, &r.DateModified)
		if err != nil {
			log.Fatal(err)
		}
		strings = append(strings, r)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	s.Logger.Logf("Fetched %d strings\n", len(strings))

	if len(strings) == 0 {
		return []core.String{}, nil
	}

	return strings, nil
}

func (s *StringRepository) CreateOne(coreString core.String) (core.String, error) {
	sql := "insert into string (name, \"order\", thread, description) " +
		"VALUES ($1, $2, $3, $4) " +
		"RETURNING *"

	var cs core.String
	err := s.DB.QueryRow(context.Background(), sql, coreString.Name, coreString.Order, coreString.Thread, coreString.Description).
		Scan(&cs.Id, &cs.Name, &cs.Order, &cs.Thread, &cs.Description, &cs.DateCreated, &cs.DateModified)

	return cs, err
}

func (s *StringRepository) DeleteById(id uuid.UUID) error {
	// TODO(Check if exists first, so you can let client know he did what was expected)
	sql := "delete from string where id = $1"
	_, err := s.DB.Exec(context.Background(), sql, id)
	return err
}

func (s *StringRepository) DeleteAllByThread(threadId uuid.UUID) error {
	sql := "delete from string where thread = $1"
	_, err := s.DB.Exec(context.Background(), sql, threadId)
	return err
}

func (s *StringRepository) UpdateName(stringId uuid.UUID, name string) error {
	sql := "update string set name = $1 where id = $2"
	res, err := s.DB.Exec(context.Background(), sql, name, stringId)
	fmt.Println("UpdateName Exec: ", res.String())

	return err
}

func (s *StringRepository) UpdateOrder(stringOrders []core.StringOrder) error {
	ctx := context.Background() // is it safe to share context with multiple functions?

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return err
	}

	// Rollback is safe to call even if the tx is already closed, so if
	// the tx commits successfully, this is a no-op
	defer tx.Rollback(ctx)

	sql := "update string set \"order\" = $1 where id = $2"

	for _, stringOrder := range stringOrders {
		_, err = tx.Exec(context.Background(), sql, stringOrder.Order, stringOrder.Id)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return err
}
