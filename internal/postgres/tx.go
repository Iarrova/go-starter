package postgres

import (
	"context"
	"errors"
	"fmt"

	ent "github.com/Iarrova/go-starter/ent/generated"
)

type txKey struct{}

// Transactor runs work in a database transaction that every repository joins
// through clientFrom. A service that needs atomicity declares its own small
// interface, e.g.
//
//	type TxRunner interface {
//		WithTx(ctx context.Context, fn func(ctx context.Context) error) error
//	}
//
// and app.NewServer passes it a *Transactor built from the same ent client.
type Transactor struct {
	client *ent.Client
}

func NewTransactor(client *ent.Client) *Transactor {
	return &Transactor{client: client}
}

// WithTx runs fn in a transaction. Repositories called with the ctx passed to fn
// take part in it. If ctx already carries a transaction, fn joins that one.
func (t *Transactor) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	if _, ok := ctx.Value(txKey{}).(*ent.Tx); ok {
		return fn(ctx)
	}

	tx, err := t.client.Tx(ctx)
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(context.WithValue(ctx, txKey{}, tx)); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			return errors.Join(err, fmt.Errorf("rolling back transaction: %w", rerr))
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	return nil
}

// clientFrom returns the transaction's client if ctx carries one, otherwise fallback.
func clientFrom(ctx context.Context, fallback *ent.Client) *ent.Client {
	if tx, ok := ctx.Value(txKey{}).(*ent.Tx); ok {
		return tx.Client()
	}
	return fallback
}
