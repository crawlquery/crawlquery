package migration

import (
	"database/sql"
	"fmt"
	"time"
)

type Migration struct {
	Name string
	SQL  string
}

var migrations = []Migration{
	{
		Name: "create_accounts_table",
		SQL: `CREATE TABLE accounts (
		id VARCHAR(36) PRIMARY KEY,
		email VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL)`,
	},
	{
		Name: "create_nodes_table",
		SQL: `CREATE TABLE nodes (
		id VARCHAR(36) PRIMARY KEY,
		account_id VARCHAR(36) NOT NULL,
		` + "`key`" + ` VARCHAR(36) NOT NULL UNIQUE,
		hostname VARCHAR(255) NOT NULL,
		port INT UNSIGNED NOT NULL,
		shard_id INT UNSIGNED NOT NULL,
		created_at TIMESTAMP NOT NULL)`,
	},
	{
		Name: "add_unique_hostname_port_index_to_nodes_table",
		SQL:  `CREATE UNIQUE INDEX nodes_account_id_hostname_port ON nodes (account_id, hostname, port)`,
	},
	{
		Name: "create_crawl_jobs_table",
		SQL: `CREATE TABLE crawl_jobs (
		id VARCHAR(36) PRIMARY KEY,
		url TEXT NOT NULL,
		page_id VARCHAR(32) NOT NULL UNIQUE,
		backoff_until TIMESTAMP,
		last_crawled_at TIMESTAMP,
		failed_reason VARCHAR(255),
		created_at TIMESTAMP NOT NULL)`,
	},
	{
		Name: "create_shards_table",
		SQL: `CREATE TABLE shards (
		id INT PRIMARY KEY,
		created_at TIMESTAMP NOT NULL)`,
	},
	{
		Name: "create_crawl_restrictions_table",
		SQL: `CREATE TABLE crawl_restrictions (
		domain VARCHAR(255) NOT NULL PRIMARY KEY,
		until TIMESTAMP)`,
	},
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

	for _, migration := range migrations {

		checkMigration := `SELECT name FROM migrations WHERE name = ?`

		row := db.QueryRow(checkMigration, migration.Name)

		var name string
		err := row.Scan(&name)

		if err != nil && err != sql.ErrNoRows {
			return err
		}

		if name == migration.Name {
			fmt.Printf("Migration %s already applied, skipping\n", migration.Name)
			continue
		}

		_, err = db.Exec(migration.SQL)
		if err != nil {
			return err
		}

		_, err = db.Exec("INSERT INTO migrations (name, created_at) VALUES (?, ?)", migration.Name, time.Now())

		if err != nil {
			return err
		}

		fmt.Printf("Applied migration %s\n", migration.Name)
	}

	return nil
}
