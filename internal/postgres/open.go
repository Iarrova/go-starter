package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib"

	ent "github.com/Iarrova/go-starter/ent/generated"
	"github.com/Iarrova/go-starter/internal/config"
)

// Open returns the ent client and its underlying *sql.DB. Closing the client
// also closes the *sql.DB.
func Open(ctx context.Context, cfg config.DatabaseConfig) (*ent.Client, *sql.DB, error) {
	db, err := sql.Open("pgx", cfg.URL)
	if err != nil {
		return nil, nil, fmt.Errorf("opening database: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		_ = db.Close() // The ping error below is the one worth reporting
		return nil, nil, fmt.Errorf("connecting to database: %w", err)
	}

	return ent.NewClient(ent.Driver(entsql.OpenDB(dialect.Postgres, db))), db, nil
}
