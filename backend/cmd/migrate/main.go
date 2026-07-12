package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5"
)

const version = 1

func main() {
	if len(os.Args) != 2 || (os.Args[1] != "up" && os.Args[1] != "down") {
		fmt.Fprintln(os.Stderr, "usage: migrate [up|down]")
		os.Exit(2)
	}
	databaseURL := os.Getenv("CROWNFALL_DATABASE_URL")
	if databaseURL == "" {
		fmt.Fprintln(os.Stderr, "CROWNFALL_DATABASE_URL is required")
		os.Exit(2)
	}
	directory := os.Getenv("CROWNFALL_MIGRATIONS_PATH")
	if directory == "" {
		directory = "migrations"
	}
	if err := run(context.Background(), databaseURL, directory, os.Args[1]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, databaseURL, directory, direction string) error {
	connection, err := pgx.Connect(ctx, databaseURL)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer connection.Close(ctx)
	if _, err := connection.Exec(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (version bigint PRIMARY KEY, applied_at timestamptz NOT NULL DEFAULT now())`); err != nil {
		return err
	}
	var exists bool
	if err := connection.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM schema_migrations WHERE version=$1)`, version).Scan(&exists); err != nil {
		return err
	}
	if direction == "up" && exists {
		return nil
	}
	if direction == "down" && !exists {
		return nil
	}
	suffix := "up"
	if direction == "down" {
		suffix = "down"
	}
	data, err := os.ReadFile(filepath.Join(directory, fmt.Sprintf("0001_createCoreTables.%s.sql", suffix)))
	if err != nil {
		return err
	}
	transaction, err := connection.Begin(ctx)
	if err != nil {
		return err
	}
	defer transaction.Rollback(ctx)
	if _, err := transaction.Exec(ctx, string(data)); err != nil {
		return fmt.Errorf("apply migration: %w", err)
	}
	if direction == "up" {
		_, err = transaction.Exec(ctx, `INSERT INTO schema_migrations(version) VALUES ($1)`, version)
	} else {
		_, err = transaction.Exec(ctx, `DELETE FROM schema_migrations WHERE version=$1`, version)
	}
	if err != nil {
		return err
	}
	return transaction.Commit(ctx)
}
