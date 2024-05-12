package migration

import (
	"database/sql"
	"time"
)

var migrations = map[string]string{
	"create_crawl_jobs_table": `CREATE TABLE crawl_jobs (
		id VARCHAR(36) PRIMARY KEY,
		url VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL
	)`,
}

var migrationTable = `CREATE TABLE IF NOT EXISTS migrations (
	id INT AUTO_INCREMENT PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	created_at TIMESTAMP NOT NULL
)`

func Up(db *sql.DB) error {
	_, err := db.Exec(migrationTable)
	if err != nil {
		return err
	}

	for k, v := range migrations {

		checkMigration := `SELECT * FROM migrations WHERE name = ?`

		row := db.QueryRow(checkMigration, k)

		var id int
		var name string
		var createdAt time.Time

		err := row.Scan(&id, &name, &createdAt)

		if err == nil {
			continue
		}

		_, err = db.Exec(v)
		if err != nil {
			return err
		}

		_, err = db.Exec("INSERT INTO migrations (name, created_at) VALUES (?, ?)", k, time.Now())

		if err != nil {
			return err
		}
	}

	return nil
}
