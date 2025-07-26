package dbsql

import (
	"fmt"
	"log"
	"os"
	"time"

	"go-phonebook/migrations"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, name)

	var db *gorm.DB
	var err error

	// Retry connection for up to 10 attempts
	for i := 1; i <= 10; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Println("Connected to MySQL database")
			break
		}
		log.Printf("Failed to connect to DB (attempt %d/10): %v", i, err)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB after retries: %w", err)
	}

	// Run migrations
	if err := migrations.MigrateTables(db); err != nil {
		log.Fatal("Migration failed:", err)
	}

	return db, nil
}
