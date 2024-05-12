package migration

import (
	"database/sql"
	"time"
)

var migrations = map[string]string{
	"create_accounts_table": `CREATE TABLE accounts (
		id VARCHAR(36) PRIMARY KEY,
		email VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL
	)`,
	"create_crawl_jobs_table": `CREATE TABLE crawl_jobs (
		id VARCHAR(36) PRIMARY KEY,
		url VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL
	)`,
	// "create_nodes_table": `CREATE TABLE nodes (
	// 	id VARCHAR(36) PRIMARY KEY,
	// 	account_id VARCHAR(36) NOT NULL,
	// 	hostname VARCHAR(255) NOT NULL,
	// 	port INT UNSIGNED NOT NULL,
	// 	created_at TIMESTAMP NOT NULL
	// )`,
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

		checkMigration := `SELECT name FROM migrations WHERE name = ?`

		row := db.QueryRow(checkMigration, k)

		var name string
		err := row.Scan(&name)

		if err != nil {
			return err
		}

		if name == k {
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
