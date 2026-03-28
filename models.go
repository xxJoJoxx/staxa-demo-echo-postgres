package main

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Contact struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func GetAllContacts(ctx context.Context, db *pgxpool.Pool) ([]Contact, error) {
	rows, err := db.Query(ctx, "SELECT id, name, email, phone, notes, created_at, updated_at FROM contacts ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []Contact
	for rows.Next() {
		var c Contact
		if err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.Notes, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		contacts = append(contacts, c)
	}
	return contacts, rows.Err()
}

func GetContactByID(ctx context.Context, db *pgxpool.Pool, id int) (*Contact, error) {
	var c Contact
	err := db.QueryRow(ctx,
		"SELECT id, name, email, phone, notes, created_at, updated_at FROM contacts WHERE id = $1", id,
	).Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.Notes, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func CreateContactDB(ctx context.Context, db *pgxpool.Pool, c *Contact) error {
	return db.QueryRow(ctx,
		"INSERT INTO contacts (name, email, phone, notes) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at",
		c.Name, c.Email, c.Phone, c.Notes,
	).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

func UpdateContactDB(ctx context.Context, db *pgxpool.Pool, c *Contact) error {
	_, err := db.Exec(ctx,
		"UPDATE contacts SET name = $1, email = $2, phone = $3, notes = $4, updated_at = NOW() WHERE id = $5",
		c.Name, c.Email, c.Phone, c.Notes, c.ID,
	)
	return err
}

func DeleteContactDB(ctx context.Context, db *pgxpool.Pool, id int) error {
	_, err := db.Exec(ctx, "DELETE FROM contacts WHERE id = $1", id)
	return err
}
