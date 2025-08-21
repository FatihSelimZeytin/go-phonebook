# Go Phonebook API

Run a secure, RESTful contact management API on your machine.

---

A full-featured phonebook backend with user auth, contact/phone management, soft deletes, Swagger docs, and more!

---

Powered by:

- Echo Web Framework
- GORM ORM
- JWT Auth
- Swagger Docs (via swaggo)
- MySQL / SQLite support

---

## Quickstart

### Docker Compose

You’ll need Docker installed.

```bash
git clone https://github.com/fatihselimzeytin/go-phonebook
cd go-phonebook
cp .env.example .env
docker compose up --build
```

---

Swagger UI will be available at: http://localhost:8080/swagger/index.html

---

## Manual Setup (Go Native)

1. Install Go
   Make sure Go ≥ 1.18 is installed:
   https://go.dev/dl/

2. Clone & Build

```bash
git clone https://github.com/fatihselimzeytin/go-phonebook
cd go-phonebook
go mod tidy
go run main.go
```

You can now visit:

* Swagger UI: http://localhost:8090/swagger/index.html
* API Root: http://localhost:8090/

---

## API Features
### Auth (JWT)
* -`/auth/register` — Register new users
* `/auth/login` — Get JWT tokens
### Contacts
* `POST /contacts` — Create contact
* `GET /contacts` — List all user contacts
* `GET /contacts/:id` — Get contact by ID
* `PUT /contacts/:id` — Update contact
* `DELETE /contacts/:id` — Soft delete contact (sets status=false)
### Phones
* Each contact can have multiple phones
* `phones` is a nested array in contact JSON

---

## Soft Delete Mechanism
Contacts aren’t deleted — they’re soft-deleted by flipping the status field to false.
```bash
Status bool `gorm:"not null;default:true" json:"status"`
```
This allows easy recovery and archiving.

---

## Migrations

Manual, versioned SQL migrations for Go Phonebook.

To run migrations:

```bash
# migrations run automatically on server start
go run main.go
```

### Adding a Migration

Create a new file in migrations/:

```go
20250821120000-create-example.go
```

```go
func DbUp0004(tx *sql.Tx) error {
    _, err := tx.Exec(`
        CREATE TABLE IF NOT EXISTS example (
            id BIGINT AUTO_INCREMENT PRIMARY KEY,
            name VARCHAR(255) NOT NULL
        );
    `)
    return err
}

func DbDown0004(tx *sql.Tx) error {
    _, err := tx.Exec(`DROP TABLE IF EXISTS example;`)
    return err
}
```

```go
var AllMigrations = []Migration{
...,
{
ID:   "20250821120000-create-example",
Up:   DbUp0004,
Down: DbDown0004,
},
}
```
The migration will now be applied automatically the next time the server starts, or when you run migrations.Migrate(db) manually.

---

## Swagger Docs
Generated using `swaggo/swag`.

To update docs:
```bash
swag init
```
Then view at: http://localhost:8090/swagger/index.html


---

## Auth Flow (JWT)
All routes under `/contacts` require authentication via `Authorization: Bearer <token>` header.

---

