package factory

import (
	"crawlquery/api/account/repository/mem"
	"crawlquery/api/account/service"
	"crawlquery/api/domain"
	"crawlquery/pkg/testutil"
	"crawlquery/pkg/util"
	"time"
)

func ValidAccount() *domain.Account {
	return &domain.Account{
		ID:        util.UUIDString(),
		Email:     "test@example.com",
		Password:  "password",
		CreatedAt: time.Now(),
	}
}

func AccountRepoWithAccount(acc *domain.Account) domain.AccountRepository {
	repo := mem.NewRepository()

	if acc != nil {
		if acc.ID == "" {
			acc.ID = util.UUIDString()
		}
		if acc.CreatedAt.IsZero() {
			acc.CreatedAt = time.Now()
		}

		if acc.Email == "" {
			acc.Email = "test@example.com"
		}

		if acc.Password == "" {
			acc.Password = "password"
		}
		repo.Create(acc)
	}

	return repo
}

func AccountServiceWithAccount(
	acc *domain.Account,
) (domain.AccountService, domain.AccountRepository) {
	repo := AccountRepoWithAccount(acc)
	return service.NewService(repo, testutil.NewTestLogger()), repo
}
