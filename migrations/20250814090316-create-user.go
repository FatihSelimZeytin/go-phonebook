package migrations

import (
	"database/sql"
)

func DbUp0001(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(255) UNIQUE NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP NULL,
			INDEX idx_users_deleted_at (deleted_at)
		);
	`)
	return err
}

func DbDown0001(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE IF EXISTS users;`)
	return err
}
