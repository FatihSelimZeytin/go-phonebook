package migrations

import (
	"database/sql"
	"errors"
	"fmt"
)

// Migration struct
type Migration struct {
	ID   string
	Up   func(tx *sql.Tx) error
	Down func(tx *sql.Tx) error
}

// AllMigrations List all migrations
var AllMigrations = []Migration{
	{
		ID:   "20250820090316-create-user",
		Up:   DbUp0001,
		Down: DbDown0001,
	},
	{
		ID:   "20250820090430-create-contacts",
		Up:   DbUp0002,
		Down: DbDown0002,
	},
	{
		ID:   "20250820090545-create-phones",
		Up:   DbUp0003,
		Down: DbDown0003,
	},
}

// Migrate runs all migrations in order, tracking applied ones
func Migrate(db *sql.DB) error {
	// Ensure the schema_migrations table exists
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			id VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	for _, m := range AllMigrations {
		// Check if this migration has already been applied
		var exists int
		err := db.QueryRow("SELECT 1 FROM schema_migrations WHERE id = ?", m.ID).Scan(&exists)
		if err == nil {
			// Already applied, skip
			fmt.Printf("Skipping already applied migration: %s\n", m.ID)
			continue
		} else if !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("failed to check migration %s: %w", m.ID, err)
		}

		// Begin transaction
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin tx for %s: %w", m.ID, err)
		}

		// Run migration
		if err := m.Up(tx); err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				return fmt.Errorf("migration %s failed: %v, rollback also failed: %w", m.ID, err, rbErr)
			}
			return fmt.Errorf("migration %s failed: %w", m.ID, err)
		}

		// Record migration as applied
		if _, err := tx.Exec("INSERT INTO schema_migrations (id) VALUES (?)", m.ID); err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				return fmt.Errorf("failed to record migration %s: %v, rollback also failed: %w", m.ID, err, rbErr)
			}
			return fmt.Errorf("failed to record migration %s: %w", m.ID, err)
		}

		// Commit transaction
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit failed for %s: %w", m.ID, err)
		}

		fmt.Printf("Migration applied: %s\n", m.ID)
	}

	return nil
}
