package repository_test

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EdlanioJ/kbu/payments/domain/entity"
	"github.com/EdlanioJ/kbu/payments/infra/db/gorm/repository"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func NewAccountTestMock() (*repository.AccountRepositoryGORM, sqlmock.Sqlmock, *entity.Account) {
	account, _ := entity.NewAccount(3400.00)

	db, mock, err := sqlmock.New()

	if err != nil {
		panic(err)
	}

	gdb, err := gorm.Open("postgres", db)

	gdb.LogMode(false)
	if err != nil {
		panic(err)
	}

	repo := repository.NewAccountRepository(gdb)

	return repo, mock, account
}

func TestAccountRepository(t *testing.T) {
	t.Parallel()

	t.Run("should test find", func(t *testing.T) {
		repo, mock, account := NewAccountTestMock()
		is := require.New(t)

		row := sqlmock.NewRows([]string{"id", "balance", "created_at"}).
			AddRow(account.ID, account.Balance, account.CreatedAt)

		const sql = `SELECT * FROM "accounts" WHERE (id = $1) ORDER BY "accounts"."id" ASC LIMIT 1`

		mock.ExpectQuery(regexp.QuoteMeta(sql)).
			WithArgs(account.ID).
			WillReturnRows(row)

		result, err := repo.Find(account.ID)

		is.Nil(err)
		is.Equal(result.ID, account.ID)
		is.Equal(result.Balance, account.Balance)
		is.Equal(result.CreatedAt, account.CreatedAt)

		id := uuid.NewV4().String()
		result, err = repo.Find(id)

		is.NotNil(err)
		is.Nil(result)
	})

	t.Run("should test save", func(t *testing.T) {
		repo, mock, account := NewAccountTestMock()
		is := require.New(t)

		rows := sqlmock.NewRows([]string{"id", "balance", "created_at"}).
			AddRow(account.ID, account.Balance, account.CreatedAt)

		const sqlSelect = `SELECT * FROM "accounts" WHERE "accounts"."id" = $1 ORDER BY "accounts"."id" ASC LIMIT 1`
		const sqlUpdate = `UPDATE "accounts" SET "created_at" = $1, "updated_at" = $2, "balance" = $3 WHERE "accounts"."id" = $4`

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(sqlUpdate)).
			WithArgs(account.CreatedAt, sqlmock.AnyArg(), account.Balance, account.ID).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		mock.ExpectQuery(regexp.QuoteMeta(sqlSelect)).
			WithArgs(account.ID).
			WillReturnRows(rows)

		err := repo.Save(account)

		is.Nil(err)

		err = repo.Save(&entity.Account{})

		is.NotNil(err)
		is.Error(err)
	})
}
