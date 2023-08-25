package sqldb

// Begin starts a new transaction.
func (s *Store) Begin() error {
	if s.Tx != nil {
		return ErrTransactionInProgress
	}

	tx, err := s.Db.Beginx()
	if err != nil {
		return err
	}

	s.Tx = tx

	return nil
}

// Commit commits the current transaction and clears the Tx object from the Store.
func (s *Store) Commit() error {
	if s.Tx == nil {
		return ErrTransactionNotStarted
	}

	defer func() {
		s.Tx = nil
	}()

	return s.Tx.Commit()
}

// Rollback calls the transaction rollback method and clears the Tx object from the Store.
func (s *Store) Rollback() error {
	if s.Tx == nil {
		return ErrTransactionNotStarted
	}

	defer func() {
		s.Tx = nil
	}()

	return s.Tx.Rollback()
}

type TestTransactor struct{}

func (t *TestTransactor) Begin() error {
	return nil
}

func (t *TestTransactor) Commit() error {
	return nil
}

func (t *TestTransactor) Rollback() error {
	return nil
}

func NewTestTransactor() *TestTransactor {
	return &TestTransactor{}
}
