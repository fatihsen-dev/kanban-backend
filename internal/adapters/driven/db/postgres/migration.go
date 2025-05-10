package db

import (
	"database/sql"
	"log"
)

func Migrate(db *sql.DB) {

	query := `CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		is_admin BOOLEAN NOT NULL DEFAULT FALSE,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	query = `CREATE TABLE IF NOT EXISTS projects (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name VARCHAR(255) NOT NULL,
		owner_id UUID NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE
	)`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	query = `CREATE TABLE IF NOT EXISTS columns (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name VARCHAR(255) NOT NULL,
		color VARCHAR(50) DEFAULT NULL,
		project_id UUID NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
	)`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	query = `CREATE TABLE IF NOT EXISTS tasks (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		title VARCHAR(255) NOT NULL,
		content TEXT DEFAULT NULL,
		column_id UUID NOT NULL,
		project_id UUID NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (column_id) REFERENCES columns(id) ON DELETE CASCADE,
		FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
	)`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	query = `CREATE TABLE IF NOT EXISTS teams (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name VARCHAR(255) NOT NULL,
		project_id UUID NOT NULL,
		role VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
	)`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	query = `CREATE TABLE IF NOT EXISTS project_members (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		project_id UUID NOT NULL,
		user_id UUID NOT NULL,
		role VARCHAR(255) NOT NULL,
		team_id UUID,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (team_id) REFERENCES teams(id)
	)`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	query = `CREATE TABLE IF NOT EXISTS invitations (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		inviter_id UUID NOT NULL,
		invitee_id UUID NOT NULL,
		project_id UUID NOT NULL,
		message VARCHAR(255) DEFAULT NULL,
		status VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (inviter_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (invitee_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
	)`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}
