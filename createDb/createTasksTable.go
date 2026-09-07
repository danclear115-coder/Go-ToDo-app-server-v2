package database

func createTasksTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			content TEXT,
			is_completed BOOLEAN NOT NULL DEFAULT 0,
			user_id INTEGER NOT NULL
		);
	`
	
	_, err := DB.Exec(query)

	return err
}