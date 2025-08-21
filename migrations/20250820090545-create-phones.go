package migrations

import "database/sql"

func DbUp0003(tx *sql.Tx) error {
	_, err := tx.Exec(`
        CREATE TABLE IF NOT EXISTS phones (
            id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
            number VARCHAR(50) NOT NULL,
            contact_id BIGINT UNSIGNED NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP,
            deleted_at TIMESTAMP NULL,
            CONSTRAINT fk_contact FOREIGN KEY (contact_id) REFERENCES contacts(id) ON DELETE CASCADE
        );
    `)
	return err
}

func DbDown0003(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE IF EXISTS phones;`)
	return err
}
