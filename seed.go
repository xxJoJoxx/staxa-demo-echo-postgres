package main

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedContacts(ctx context.Context, db *pgxpool.Pool) error {
	var count int
	err := db.QueryRow(ctx, "SELECT COUNT(*) FROM contacts").Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	seeds := []Contact{
		{Name: "Alice Johnson", Email: "alice@example.com", Phone: "555-0101", Notes: "Met at the Go meetup"},
		{Name: "Bob Smith", Email: "bob@example.com", Phone: "555-0102", Notes: "Colleague from previous job"},
		{Name: "Carol Williams", Email: "carol@example.com", Phone: "555-0103", Notes: "Friend from college"},
	}

	for _, c := range seeds {
		if err := CreateContactDB(ctx, db, &c); err != nil {
			return err
		}
	}
	return nil
}
