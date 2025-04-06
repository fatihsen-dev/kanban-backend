package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	DB *sql.DB
}

func NewPostgresRepository(connStr string) *PostgresRepository {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	Migrate(db)

	return &PostgresRepository{DB: db}
}

func (r *PostgresRepository) Close() error {
	return r.DB.Close()
}
