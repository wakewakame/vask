package db

import (
	"database/sql"
	"fmt"

	"vask/internal/model"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	db *sql.DB
}

func Open(path string) (*DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("Failed to open database: %w", err)
	}

	_, err = db.Exec(`
		PRAGMA foreign_keys = ON;

		CREATE TABLE IF NOT EXISTS projects (
			id         INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			name       TEXT    NOT NULL UNIQUE,
			created_at TEXT    NOT NULL DEFAULT (DATETIME('now')),
			updated_at TEXT    NOT NULL DEFAULT (DATETIME('now'))
		);
		CREATE TRIGGER IF NOT EXISTS trigger_projects_updated_at AFTER UPDATE ON projects
		BEGIN
			UPDATE projects SET updated_at = DATETIME('now') WHERE rowid == NEW.rowid;
		END;

		CREATE TABLE IF NOT EXISTS versions (
			id                    INTEGER    NOT NULL PRIMARY KEY AUTOINCREMENT,
			project               INTEGER,
			content               TEXT       NOT NULL,
			created_at            TEXT       NOT NULL DEFAULT (DATETIME('now')),
			updated_at            TEXT       NOT NULL DEFAULT (DATETIME('now')),
			FOREIGN KEY (project) REFERENCES projects(id)
		);
		CREATE TRIGGER IF NOT EXISTS trigger_versions_updated_at AFTER UPDATE ON versions
		BEGIN
			UPDATE projects SET updated_at = DATETIME('now') WHERE rowid == NEW.rowid;
		END;
	`)
	if err != nil {
		return nil, fmt.Errorf("Failed to create tables: %w", err)
	}

	return &DB{db: db}, nil
}

func (db *DB) Close() {
	defer db.db.Close()
}

func (db *DB) GetProjects() ([]*model.Project, error) {
	rows, err := db.db.Query("SELECT id, name, created_at, updated_at FROM projects")
	if err != nil {
		return nil, fmt.Errorf("Failed to exec query: %w", err)
	}
	defer rows.Close()

	result := make([]*model.Project, 0)
	for rows.Next() {
		project := model.Project{}
		var createdAtStr string
		var updatedAtStr string
		err = rows.Scan(&project.Id, &project.Name, &createdAtStr, &updatedAtStr)
		if err != nil {
			return nil, fmt.Errorf("Failed to scan row: %w", err)
		}

		project.CreatedAt, project.UpdatedAt, err = stringToTime(createdAtStr, updatedAtStr)
		if err != nil {
			return nil, err
		}
		result = append(result, &project)
	}

	return result, nil
}

func (db *DB) AddProject(name string) (*int64, error) {
	stmt, err := db.db.Prepare("INSERT INTO projects(name) values(?)")
	if err != nil {
		return nil, fmt.Errorf("Failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(name)
	if err != nil {
		return nil, fmt.Errorf("Failed to exec: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("Failed to get id: %w", err)
	}

	return &id, nil
}

func (db *DB) GetProject(id int64) (*model.Project, error) {
	stmt, err := db.db.Prepare("SELECT name, created_at, updated_at FROM projects WHERE id = ?")
	if err != nil {
		return nil, fmt.Errorf("Failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	project := model.Project{Id: id}
	var createdAtStr string
	var updatedAtStr string
	err = stmt.QueryRow(id).Scan(&project.Name, &createdAtStr, &updatedAtStr)
	project.CreatedAt, project.UpdatedAt, err = stringToTime(createdAtStr, updatedAtStr)

	if err != nil {
		return nil, fmt.Errorf("Failed to exec: %w", err)
	}

	return &project, nil
}

func (db *DB) SetProject(id int64, name string) error {
	stmt, err := db.db.Prepare("UPDATE projects SET name = ? WHERE id = ?")
	if err != nil {
		return fmt.Errorf("Failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, id)
	if err != nil {
		return fmt.Errorf("Failed to exec: %w", err)
	}

	return nil
}

func (db *DB) DeleteProject(id int64) error {
	stmt, err := db.db.Prepare("DELETE FROM projects WHERE id = ?")
	if err != nil {
		return fmt.Errorf("Failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("Failed to exec: %w", err)
	}

	return nil
}

func stringToTime(createdAtStr, updatedAtStr string) (model.Time, model.Time, error) {
	const layout = "2006-01-02 15:04:05"
	createdAt, err := model.ParseTime(layout, createdAtStr)
	zero := model.Time{}
	if err != nil {
		return zero, zero, fmt.Errorf("Failed to parse created time (%s): %w", createdAtStr, err)
	}
	updatedAt, err := model.ParseTime(layout, updatedAtStr)
	if err != nil {
		return zero, zero, fmt.Errorf("Failed to parse updated time (%s): %w", updatedAtStr, err)
	}
	return createdAt, updatedAt, nil
}
