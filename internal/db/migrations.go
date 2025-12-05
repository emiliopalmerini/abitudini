package db

import (
	"database/sql"
	"fmt"
)

func Migrate(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS habits (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		description TEXT NOT NULL,
		start_date TEXT NOT NULL,
		color TEXT NOT NULL,
		created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS records (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		habit_id INTEGER NOT NULL,
		record_date TEXT NOT NULL,
		completed_at TEXT NOT NULL,
		created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (habit_id) REFERENCES habits(id) ON DELETE CASCADE,
		UNIQUE(habit_id, record_date)
	);

	CREATE INDEX IF NOT EXISTS idx_records_habit_id ON records(habit_id);
	CREATE INDEX IF NOT EXISTS idx_records_date ON records(record_date);
	`

	if _, err := db.Exec(schema); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	return nil
}
