package main

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RunMigrations(ctx context.Context, db *pgxpool.Pool) error {
	query := `
	CREATE TABLE IF NOT EXISTS contacts (
		id          SERIAL PRIMARY KEY,
		name        VARCHAR(255) NOT NULL,
		email       VARCHAR(255) NOT NULL,
		phone       VARCHAR(50),
		notes       TEXT,
		created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);`
	_, err := db.Exec(ctx, query)
	return err
}
