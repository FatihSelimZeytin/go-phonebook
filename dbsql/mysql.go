package dbsql

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"go-phonebook/migrations"

	_ "github.com/go-sql-driver/mysql" // Import MySQL driver
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ConnectGorm returns a *gorm.DB for app usage
func ConnectGorm() (*gorm.DB, error) {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, name)

	var db *gorm.DB
	var err error

	for i := 1; i <= 10; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Println("Connected to MySQL via GORM")
			break
		}
		log.Printf("Failed to connect to DB (attempt %d/10): %v", i, err)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB after retries: %w", err)
	}

	return db, nil
}

// ConnectSQL returns a *sql.DB for migrations
func ConnectSQL() (*sql.DB, error) {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, name)

	var db *sql.DB
	var err error

	for i := 1; i <= 10; i++ {
		db, err = sql.Open("mysql", dsn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				log.Println("Connected to MySQL via database/sql")
				break
			}
		}
		log.Printf("Failed to connect to DB (attempt %d/10): %v", i, err)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB after retries: %w", err)
	}

	return db, nil
}

// RunMigrations runs all custom migrations using *sql.DB
func RunMigrations() error {
	db, err := ConnectSQL()
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	return migrations.Migrate(db)
}
