package db

import (
	"database/sql"
	"fmt"
)

const (
	FrequencyDaily   = "daily"
	FrequencyWeekly  = "weekly"
	FrequencyMonthly = "monthly"
)

func Migrate(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS habits (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		description TEXT NOT NULL,
		frequency TEXT NOT NULL CHECK(frequency IN ('daily', 'weekly', 'monthly')),
		start_date TEXT NOT NULL,
		color TEXT NOT NULL,
		created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS habit_schedule (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		habit_id INTEGER NOT NULL,
		day_of_week INTEGER,
		day_of_month INTEGER,
		created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (habit_id) REFERENCES habits(id) ON DELETE CASCADE,
		UNIQUE(habit_id, day_of_week, day_of_month)
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
	CREATE INDEX IF NOT EXISTS idx_habit_schedule_habit_id ON habit_schedule(habit_id);
	`

	if _, err := db.Exec(schema); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	return nil
}
