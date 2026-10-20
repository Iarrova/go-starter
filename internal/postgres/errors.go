package postgres

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// SQLSTATE codes, see https://www.postgresql.org/docs/current/errcodes-appendix.html
const (
	foreignKeyViolation = "23503"
	uniqueViolation     = "23505"
)

// isConstraintViolation reports whether err is a Postgres error with the given
// SQLSTATE code on the named constraint. ent wraps driver errors, so this
// unwraps to the *pgconn.PgError underneath.
func isConstraintViolation(err error, code, constraint string) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == code && pgErr.ConstraintName == constraint
}
	