package db

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

const pgUniqueViolation = "23505"

type UniqueViolation struct {
	Constraint string
}

func AsUniqueViolation(err error) (*UniqueViolation, bool) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
		return &UniqueViolation{Constraint: pgErr.ConstraintName}, true
	}
	return nil, false
}
