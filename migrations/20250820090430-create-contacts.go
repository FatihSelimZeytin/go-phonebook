package migrations

import (
	"database/sql"
)

func DbUp0002(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS contacts (
			id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
			first_name VARCHAR(255) NOT NULL,
			surname VARCHAR(255) NOT NULL,
			company VARCHAR(255) NULL,
			user_id BIGINT UNSIGNED NOT NULL,
			status BOOLEAN NOT NULL DEFAULT true,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP,
			CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);
	`)
	return err
}

func DbDown0002(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE IF EXISTS contacts;`)
	return err
}
