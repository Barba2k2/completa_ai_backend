package database

import (
	"database/sql"
	"embed"
	"log"
	"sort"
	"strings"
)

type Migrator struct {
	db         *sql.DB
	migrations embed.FS
}

func NewMigrator(db *sql.DB, migrations embed.FS) *Migrator {
	return &Migrator{db: db, migrations: migrations}
}

func (m *Migrator) Run() error {
	if err := m.createMigrationsTable(); err != nil {
		return err
	}

	applied, err := m.getAppliedMigrations()
	if err != nil {
		return err
	}

	files, err := m.getMigrationFiles()
	if err != nil {
		return err
	}

	for _, file := range files {
		if applied[file] {
			continue
		}

		log.Printf("[Migrator] Applying migration: %s", file)

		content, err := m.migrations.ReadFile("sql/" + file)
		if err != nil {
			return err
		}

		if _, err := m.db.Exec(string(content)); err != nil {
			log.Printf("[Migrator] Failed to apply %s: %v", file, err)
			return err
		}

		if err := m.markAsApplied(file); err != nil {
			return err
		}

		log.Printf("[Migrator] Applied migration: %s", file)
	}

	return nil
}

func (m *Migrator) createMigrationsTable() error {
	_, err := m.db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMPTZ DEFAULT NOW()
		)
	`)
	return err
}

func (m *Migrator) getAppliedMigrations() (map[string]bool, error) {
	applied := make(map[string]bool)

	rows, err := m.db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}

	return applied, rows.Err()
}

func (m *Migrator) getMigrationFiles() ([]string, error) {
	var files []string

	entries, err := m.migrations.ReadDir("sql")
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			files = append(files, entry.Name())
		}
	}

	sort.Strings(files)
	return files, nil
}

func (m *Migrator) markAsApplied(version string) error {
	_, err := m.db.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", version)
	return err
}
