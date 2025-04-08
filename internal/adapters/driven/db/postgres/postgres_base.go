package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	DB *sql.DB
}

func NewPostgresRepository(connStr string) *PostgresRepository {
	defaultConnStr := strings.Replace(connStr, "/kanban?", "/postgres?", 1)
	defaultDB, err := sql.Open("postgres", defaultConnStr)
	if err != nil {
		log.Fatal(err)
	}
	defer defaultDB.Close()

	dbName := "kanban"
	if idx := strings.LastIndex(connStr, "/"); idx != -1 {
		if endIdx := strings.Index(connStr[idx+1:], "?"); endIdx != -1 {
			dbName = connStr[idx+1 : idx+1+endIdx]
		}
	}

	var exists bool
	err = defaultDB.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", dbName).Scan(&exists)
	if err != nil {
		log.Fatal(err)
	}

	if !exists {
		_, err = defaultDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			log.Fatal(err)
		}
	}

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
