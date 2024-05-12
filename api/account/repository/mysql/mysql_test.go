package mysql_test

import (
	"crawlquery/api/account/repository/mysql"
	"crawlquery/api/domain"
	"crawlquery/api/migration"
	"crawlquery/pkg/testutil"
	"crawlquery/pkg/util"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	t.Run("can create an account", func(t *testing.T) {
		// Arrange
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)

		repo := mysql.NewRepository(db)

		// Act
		account := &domain.Account{
			ID:        util.UUID(),
			Email:     "test@example.com",
			Password:  "password",
			CreatedAt: time.Now().UTC(),
		}

		defer db.Exec("DELETE FROM accounts WHERE id = ?", account.ID)

		err := repo.Create(account)

		// Assert
		if err != nil {
			t.Errorf("Error creating account: %v", err)
		}

		res, err := db.Query("SELECT * FROM accounts WHERE id = ?", account.ID)

		if err != nil {
			t.Errorf("Error querying for account: %v", err)
		}

		var id string
		var email string
		var password string
		var createdAt time.Time

		for res.Next() {
			err = res.Scan(&id, &email, &password, &createdAt)
			if err != nil {
				t.Errorf("Error scanning account: %v", err)
			}
		}

		if id != account.ID {
			t.Errorf("Expected ID to be %s, got %s", account.ID, id)
		}

		if email != account.Email {
			t.Errorf("Expected Email to be %s, got %s", account.Email, email)
		}

		if password != account.Password {
			t.Errorf("Expected Password to be %s, got %s", account.Password, password)
		}

		if createdAt.Sub(account.CreatedAt) > time.Second || account.CreatedAt.Sub(createdAt) > time.Second {
			t.Errorf("Expected CreatedAt to be within one second of %v, got %v", account.CreatedAt, createdAt)
		}

	})

	t.Run("cant create an account with the same ID", func(t *testing.T) {
		// Arrange
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)

		repo := mysql.NewRepository(db)

		account := &domain.Account{
			ID:        util.UUID(),
			Email:     "test@example.com",
			Password:  "password",
			CreatedAt: time.Now().UTC(),
		}

		err := repo.Create(account)
		defer db.Exec("DELETE FROM accounts WHERE id = ?", account.ID)

		if err != nil {
			t.Fatalf("Error creating account: %v", err)
		}

		// Act
		err = repo.Create(account)

		// Assert
		if err == nil {
			t.Errorf("Expected error creating account, got nil")
		}
	})

	t.Run("cant create an account with the same email", func(t *testing.T) {
		// Arrange
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)

		repo := mysql.NewRepository(db)

		account := &domain.Account{
			ID:        util.UUID(),
			Email:     "test@example.com",
			Password:  "password",
			CreatedAt: time.Now().UTC(),
		}

		defer db.Exec("DELETE FROM accounts WHERE id = ?", account.ID)

		err := repo.Create(account)

		if err != nil {
			t.Fatalf("Error creating account: %v", err)
		}

		// Act
		account.ID = util.UUID()

		err = repo.Create(account)

		// Assert
		if err == nil {
			t.Errorf("Expected error creating account, got nil")
		}
	})
}

func TestGet(t *testing.T) {
	t.Run("can get an account", func(t *testing.T) {
		// Arrange
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)

		repo := mysql.NewRepository(db)

		account := &domain.Account{
			ID:        util.UUID(),
			Email:     "test@example.com",
			Password:  "password",
			CreatedAt: time.Now().UTC(),
		}

		err := repo.Create(account)

		if err != nil {
			t.Fatalf("Error creating account: %v", err)
		}

		defer db.Exec("DELETE FROM accounts WHERE id = ?", account.ID)

		// Act
		res, err := repo.Get(account.ID)

		// Assert
		if err != nil {
			t.Errorf("Error getting account: %v", err)
		}

		if res.ID != account.ID {
			t.Errorf("Expected ID to be %s, got %s", account.ID, res.ID)
		}

		if res.Email != account.Email {
			t.Errorf("Expected Email to be %s, got %s", account.Email, res.Email)
		}

		if res.Password != account.Password {
			t.Errorf("Expected Password to be %s, got %s", account.Password, res.Password)
		}

		if res.CreatedAt.Sub(account.CreatedAt) > time.Second || account.CreatedAt.Sub(res.CreatedAt) > time.Second {
			t.Errorf("Expected CreatedAt to be within one second of %v, got %v", account.CreatedAt, res.CreatedAt)
		}
	})

	t.Run("cant get an account that doesnt exist", func(t *testing.T) {
		// Arrange
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)

		repo := mysql.NewRepository(db)

		// Act
		_, err := repo.Get(util.UUID())

		// Assert
		if err == nil {
			t.Errorf("Expected error getting account, got nil")
		}
	})
}

func TestGetByEmail(t *testing.T) {
	t.Run("can get an account by email", func(t *testing.T) {
		// Arrange
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)

		repo := mysql.NewRepository(db)

		account := &domain.Account{
			ID:        util.UUID(),
			Email:     "test@example.com",
			Password:  "password",
			CreatedAt: time.Now().UTC(),
		}

		err := repo.Create(account)

		if err != nil {
			t.Fatalf("Error creating account: %v", err)
		}

		defer db.Exec("DELETE FROM accounts WHERE id = ?", account.ID)

		// Act
		res, err := repo.GetByEmail(account.Email)

		// Assert
		if err != nil {
			t.Errorf("Error getting account: %v", err)
		}

		if res.ID != account.ID {
			t.Errorf("Expected ID to be %s, got %s", account.ID, res.ID)
		}

		if res.Email != account.Email {
			t.Errorf("Expected Email to be %s, got %s", account.Email, res.Email)
		}

		if res.Password != account.Password {
			t.Errorf("Expected Password to be %s, got %s", account.Password, res.Password)
		}

		if res.CreatedAt.Sub(account.CreatedAt) > time.Second || account.CreatedAt.Sub(res.CreatedAt) > time.Second {
			t.Errorf("Expected CreatedAt to be within one second of %v, got %v", account.CreatedAt, res.CreatedAt)
		}
	})

	t.Run("cant get an account by email that doesnt exist", func(t *testing.T) {
		// Arrange
		db := testutil.CreateTestMysqlDB()
		defer db.Close()
		migration.Up(db)

		repo := mysql.NewRepository(db)

		// Act
		_, err := repo.GetByEmail("test@example.com")

		// Assert
		if err == nil {
			t.Errorf("Expected error getting account, got nil")
		}
	})
}
