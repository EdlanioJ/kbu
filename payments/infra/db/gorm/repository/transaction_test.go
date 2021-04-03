package repository_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EdlanioJ/kbu/payments/domain/entity"
	"github.com/EdlanioJ/kbu/payments/infra/db/gorm/repository"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func NewTransactionTestMock() (*repository.TransactionRepositoryGORM, sqlmock.Sqlmock, *entity.Transaction) {
	transactionType := entity.TransactionToService
	accountFrom, _ := entity.NewAccount(3000)
	accountTo, _ := entity.NewAccount(200)
	externalID := uuid.NewV4().String()
	currency := "AOA"
	amount := 30.00
	transaction, _ := entity.NewTransaction(accountFrom, accountTo, externalID, transactionType, currency, amount)

	db, mock, err := sqlmock.New()

	if err != nil {
		panic(err)
	}

	gdb, err := gorm.Open("postgres", db)

	gdb.LogMode(false)
	if err != nil {
		panic(err)
	}

	repo := repository.NewTransactionRepository(gdb)

	return repo, mock, transaction
}

func TestTransactionRepository(t *testing.T) {
	t.Parallel()

	t.Run("should test regiter", func(t *testing.T) {
		repo, mock, transaction := NewTransactionTestMock()
		is := require.New(t)

		const insertSql = `INSERT INTO "transactions" ("id","created_at","updated_at","amount","status","currency","account_from_id","account_to_id","type","external_id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING "transactions"."id"`
		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(insertSql)).
			WithArgs(
				transaction.ID, transaction.CreatedAt, sqlmock.AnyArg(), transaction.Amount, transaction.Status, transaction.Currency, transaction.AccountFromID, transaction.AccountToID, transaction.Type, transaction.ExternalID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(transaction.ID))
		mock.ExpectCommit()

		err := repo.Register(transaction)
		is.Nil(err)

		err = repo.Register(&entity.Transaction{})
		is.NotNil(err)
	})

	t.Run("should test save", func(t *testing.T) {
		repo, mock, transaction := NewTransactionTestMock()
		is := require.New(t)

		row := sqlmock.NewRows([]string{"id", "account_from_id", "amount", "status", "currency", "account_to_id", "created_at", "updated_at"}).
			AddRow(transaction.ID, transaction.AccountFromID, transaction.Amount, transaction.Status, transaction.Currency, transaction.AccountToID, transaction.CreatedAt, transaction.UpdatedAt)

		const updateSql = `UPDATE "transactions" SET "created_at" = $1, "updated_at" = $2, "amount" = $3, "status" = $4, "currency" = $5, "account_from_id" = $6, "account_to_id" = $7, "type" = $8, "external_id" = $9 WHERE "transactions"."id" = $10`
		const selectTransaction = `SELECT * FROM "transactions"  WHERE "transactions"."id" = $1 ORDER BY "transactions"."id" ASC LIMIT 1`

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(updateSql)).
			WithArgs(transaction.CreatedAt, sqlmock.AnyArg(), transaction.Amount, transaction.Status, transaction.Currency, transaction.AccountFromID, transaction.AccountToID, transaction.Type, transaction.ExternalID, transaction.ID).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		mock.ExpectQuery(regexp.QuoteMeta(selectTransaction)).
			WithArgs(transaction.ID).
			WillReturnRows(row)

		err := repo.Save(transaction)

		is.Nil(err)

		err = repo.Save(&entity.Transaction{})

		is.NotNil(err)
	})

	t.Run("should test find", func(t *testing.T) {
		repo, mock, transaction := NewTransactionTestMock()
		is := require.New(t)

		row := sqlmock.NewRows([]string{"id", "account_from_id", "amount", "status", "currency", "account_to_id", "created_at", "updated_at"}).
			AddRow(transaction.ID, transaction.AccountFromID, transaction.Amount, transaction.Status, transaction.Currency, transaction.AccountToID, transaction.CreatedAt, transaction.UpdatedAt)

		const selectTransaction = `SELECT * FROM "transactions" WHERE (id = $1) ORDER BY "transactions"."id" ASC LIMIT 1`

		mock.ExpectQuery(regexp.QuoteMeta(selectTransaction)).
			WithArgs(transaction.ID).
			WillReturnRows(row)

		result, err := repo.Find(transaction.ID)

		is.Nil(err)
		is.Equal(result.ID, transaction.ID)
		is.Equal(result.AccountFromID, transaction.AccountFromID)
		is.Equal(result.Amount, transaction.Amount)

		result, err = repo.Find(uuid.NewV4().String())

		is.NotNil(err)
		is.Nil(result)
	})

	t.Run("should test find all", func(t *testing.T) {
		repo, mock, transaction := NewTransactionTestMock()
		is := require.New(t)

		row := sqlmock.NewRows([]string{"id", "account_from_id", "amount", "status", "currency", "account_to_id", "created_at", "updated_at"}).
			AddRow(transaction.ID, transaction.AccountFromID, transaction.Amount, transaction.Status, transaction.Currency, transaction.AccountToID, transaction.CreatedAt, transaction.UpdatedAt)

		countRow := sqlmock.NewRows([]string{"count"}).AddRow(1)

		page := 1
		limit := 10
		sort := "created_at DESC"

		pagination := &entity.Pagination{
			Page:  page,
			Limit: limit,
			Sort:  sort,
		}

		selectTransaction := fmt.Sprintf(`SELECT * FROM "transactions" ORDER BY %s LIMIT %d OFFSET %d`, sort, limit, 0)
		const countSelect = `SELECT count(*) FROM "transactions"`

		mock.ExpectQuery(regexp.QuoteMeta(selectTransaction)).WillReturnRows(row)
		mock.ExpectQuery(regexp.QuoteMeta(countSelect)).WillReturnRows(countRow)

		result, total, err := repo.FindAll(pagination)

		is.Nil(err)
		is.Equal(total, 1)
		is.Equal(result[0].ID, transaction.ID)
		is.Equal(result[0].AccountFromID, transaction.AccountFromID)
		is.Equal(result[0].Amount, transaction.Amount)

		result, total, err = repo.FindAll(&entity.Pagination{
			Page:  1,
			Limit: 20,
			Sort:  "created_at ASC",
		})

		is.Nil(result)
		is.Equal(total, 0)
		is.NotNil(err)
	})

	t.Run("should test find by type", func(t *testing.T) {
		repo, mock, transaction := NewTransactionTestMock()
		is := require.New(t)

		row := sqlmock.NewRows([]string{"id", "account_from_id", "amount", "status", "currency", "account_to_id", "created_at", "updated_at", "type", "external_id"}).
			AddRow(transaction.ID, transaction.AccountFromID, transaction.Amount, transaction.Status, transaction.Currency, transaction.AccountToID, transaction.CreatedAt, transaction.UpdatedAt, transaction.Type, transaction.ExternalID)

		selectTransaction := `SELECT * FROM "transactions" WHERE (id = $1 AND type = $2) ORDER BY "transactions"."id" ASC LIMIT 1`

		mock.ExpectQuery(regexp.QuoteMeta(selectTransaction)).
			WithArgs(transaction.ID, transaction.Type).
			WillReturnRows(row)

		result, err := repo.FindByType(transaction.ID, transaction.Type)

		is.Nil(err)
		is.NotNil(result)
		is.Equal(result.ID, transaction.ID)
		is.Equal(result.Type, transaction.Type)

		result, err = repo.FindByType(transaction.AccountFromID, transaction.Type)
		is.Nil(result)
		is.NotNil(err)
		is.Error(err)
	})

	t.Run("should test find all by type", func(t *testing.T) {
		repo, mock, transaction := NewTransactionTestMock()
		is := require.New(t)

		row := sqlmock.NewRows([]string{"id", "account_from_id", "amount", "status", "currency", "account_to_id", "created_at", "updated_at", "type", "external_id"}).
			AddRow(transaction.ID, transaction.AccountFromID, transaction.Amount, transaction.Status, transaction.Currency, transaction.AccountToID, transaction.CreatedAt, transaction.UpdatedAt, transaction.Type, transaction.ExternalID)

		page := 1
		limit := 10
		sort := "created_at DESC"

		pagination := &entity.Pagination{
			Page:  page,
			Limit: limit,
			Sort:  sort,
		}

		selectTransaction := fmt.Sprintf(`SELECT * FROM "transactions" WHERE (type = $1) ORDER BY %s LIMIT %d OFFSET %d`, sort, limit, 0)
		countSelect := `SELECT count(*) FROM "transactions" WHERE (type = $1)`

		countRow := sqlmock.NewRows([]string{"count"}).AddRow(1)

		mock.ExpectQuery(regexp.QuoteMeta(selectTransaction)).WithArgs(transaction.Type).WillReturnRows(row)
		mock.ExpectQuery(regexp.QuoteMeta(countSelect)).WithArgs(transaction.Type).WillReturnRows(countRow)

		result, total, err := repo.FindAllByType(transaction.Type, pagination)

		fmt.Println(err)
		is.Nil(err)
		is.Equal(total, 1)
		is.Equal(result[0].ID, transaction.ID)
		is.Equal(result[0].AccountFromID, transaction.AccountFromID)
		is.Equal(result[0].Amount, transaction.Amount)

		transactionType := entity.TransactionToService
		result, total, err = repo.FindAllByType(transactionType, &entity.Pagination{
			Page:  2,
			Limit: limit,
			Sort:  "id DESC",
		})

		is.Nil(result)
		is.Equal(total, 0)
		is.NotNil(err)
	})

	t.Run("should test find by external id", func(t *testing.T) {
		repo, mock, transaction := NewTransactionTestMock()
		is := require.New(t)

		row := sqlmock.NewRows([]string{"id", "account_from_id", "amount", "status", "currency", "account_to_id", "created_at", "updated_at", "type", "external_id"}).
			AddRow(transaction.ID, transaction.AccountFromID, transaction.Amount, transaction.Status, transaction.Currency, transaction.AccountToID, transaction.CreatedAt, transaction.UpdatedAt, transaction.Type, transaction.ExternalID)

		selectTransaction := `SELECT * FROM "transactions" WHERE (id = $1 AND external_id = $2) ORDER BY "transactions"."id" ASC LIMIT 1`

		mock.ExpectQuery(regexp.QuoteMeta(selectTransaction)).
			WithArgs(transaction.ID, transaction.ExternalID).
			WillReturnRows(row)

		result, err := repo.FindByExternalID(transaction.ID, transaction.ExternalID)

		is.Nil(err)
		is.NotNil(result)
		is.Equal(result.ID, transaction.ID)
		is.Equal(result.ExternalID, transaction.ExternalID)

		result, err = repo.FindByExternalID(transaction.AccountFromID, transaction.ExternalID)
		is.Nil(result)
		is.NotNil(err)
		is.Error(err)
	})

	t.Run("should test find all by external id", func(t *testing.T) {
		repo, mock, transaction := NewTransactionTestMock()
		is := require.New(t)

		row := sqlmock.NewRows([]string{"id", "account_from_id", "amount", "status", "currency", "account_to_id", "created_at", "updated_at", "type", "external_id"}).
			AddRow(transaction.ID, transaction.AccountFromID, transaction.Amount, transaction.Status, transaction.Currency, transaction.AccountToID, transaction.CreatedAt, transaction.UpdatedAt, transaction.Type, transaction.ExternalID)

		page := 1
		limit := 10
		sort := "created_at DESC"

		pagination := &entity.Pagination{
			Page:  page,
			Limit: limit,
			Sort:  sort,
		}

		selectTransaction := fmt.Sprintf(`SELECT * FROM "transactions" WHERE (external_id = $1) ORDER BY %s LIMIT %d OFFSET %d`, sort, limit, 0)
		countSelect := `SELECT count(*) FROM "transactions" WHERE (external_id = $1)`

		countRow := sqlmock.NewRows([]string{"count"}).AddRow(1)

		mock.ExpectQuery(regexp.QuoteMeta(selectTransaction)).WithArgs(transaction.ExternalID).WillReturnRows(row)
		mock.ExpectQuery(regexp.QuoteMeta(countSelect)).WithArgs(transaction.ExternalID).WillReturnRows(countRow)

		result, total, err := repo.FindAllByExternalID(transaction.ExternalID, pagination)

		fmt.Println(err)
		is.Nil(err)
		is.Equal(total, 1)
		is.Equal(result[0].ID, transaction.ID)
		is.Equal(result[0].AccountFromID, transaction.AccountFromID)
		is.Equal(result[0].Amount, transaction.Amount)

		result, total, err = repo.FindAllByExternalID(transaction.AccountFromID, &entity.Pagination{
			Page:  2,
			Limit: limit,
			Sort:  "id DESC",
		})

		is.Nil(result)
		is.Equal(total, 0)
		is.NotNil(err)
	})

	t.Run("should test find by account from id", func(t *testing.T) {
		repo, mock, transaction := NewTransactionTestMock()
		is := require.New(t)

		row := sqlmock.NewRows([]string{"id", "account_from_id", "amount", "status", "currency", "account_to_id", "created_at", "updated_at", "type", "external_id"}).
			AddRow(transaction.ID, transaction.AccountFromID, transaction.Amount, transaction.Status, transaction.Currency, transaction.AccountToID, transaction.CreatedAt, transaction.UpdatedAt, transaction.Type, transaction.ExternalID)

		selectTransaction := `SELECT * FROM "transactions" WHERE (id = $1 AND account_from_id = $2) ORDER BY "transactions"."id" ASC LIMIT 1`

		mock.ExpectQuery(regexp.QuoteMeta(selectTransaction)).
			WithArgs(transaction.ID, transaction.AccountFromID).
			WillReturnRows(row)

		result, err := repo.FindByFromAccountID(transaction.ID, transaction.AccountFromID)

		is.Nil(err)
		is.NotNil(result)
		is.Equal(result.ID, transaction.ID)
		is.Equal(result.AccountFromID, transaction.AccountFromID)

		result, err = repo.FindByFromAccountID(transaction.AccountFromID, transaction.AccountToID)
		is.Nil(result)
		is.NotNil(err)
		is.Error(err)
	})

	t.Run("should test find all by account from id", func(t *testing.T) {
		repo, mock, transaction := NewTransactionTestMock()
		is := require.New(t)

		row := sqlmock.NewRows([]string{"id", "account_from_id", "amount", "status", "currency", "account_to_id", "created_at", "updated_at", "type", "external_id"}).
			AddRow(transaction.ID, transaction.AccountFromID, transaction.Amount, transaction.Status, transaction.Currency, transaction.AccountToID, transaction.CreatedAt, transaction.UpdatedAt, transaction.Type, transaction.ExternalID)

		page := 1
		limit := 10
		sort := "created_at DESC"

		pagination := &entity.Pagination{
			Page:  page,
			Limit: limit,
			Sort:  sort,
		}

		selectTransaction := fmt.Sprintf(`SELECT * FROM "transactions" WHERE (account_from_id = $1) ORDER BY %s LIMIT %d OFFSET %d`, sort, limit, 0)
		countSelect := `SELECT count(*) FROM "transactions" WHERE (account_from_id = $1)`

		countRow := sqlmock.NewRows([]string{"count"}).AddRow(1)

		mock.ExpectQuery(regexp.QuoteMeta(selectTransaction)).WithArgs(transaction.AccountFromID).WillReturnRows(row)
		mock.ExpectQuery(regexp.QuoteMeta(countSelect)).WithArgs(transaction.AccountFromID).WillReturnRows(countRow)

		result, total, err := repo.FindAllByFromAccountID(transaction.AccountFromID, pagination)

		fmt.Println(err)
		is.Nil(err)
		is.Equal(total, 1)
		is.Equal(result[0].ID, transaction.ID)
		is.Equal(result[0].AccountFromID, transaction.AccountFromID)
		is.Equal(result[0].Amount, transaction.Amount)

		result, total, err = repo.FindAllByFromAccountID(transaction.ExternalID, &entity.Pagination{
			Page:  2,
			Limit: limit,
			Sort:  "id DESC",
		})

		is.Nil(result)
		is.Equal(total, 0)
		is.NotNil(err)
	})

	t.Run("should test find by account to id", func(t *testing.T) {
		repo, mock, transaction := NewTransactionTestMock()
		is := require.New(t)

		row := sqlmock.NewRows([]string{"id", "account_from_id", "amount", "status", "currency", "account_to_id", "created_at", "updated_at", "type", "external_id"}).
			AddRow(transaction.ID, transaction.AccountFromID, transaction.Amount, transaction.Status, transaction.Currency, transaction.AccountToID, transaction.CreatedAt, transaction.UpdatedAt, transaction.Type, transaction.ExternalID)

		selectTransaction := `SELECT * FROM "transactions" WHERE (id = $1 AND account_to_id = $2) ORDER BY "transactions"."id" ASC LIMIT 1`

		mock.ExpectQuery(regexp.QuoteMeta(selectTransaction)).
			WithArgs(transaction.ID, transaction.AccountToID).
			WillReturnRows(row)

		result, err := repo.FindByToAccountID(transaction.ID, transaction.AccountToID)

		is.Nil(err)
		is.NotNil(result)
		is.Equal(result.ID, transaction.ID)
		is.Equal(result.AccountFromID, transaction.AccountFromID)

		result, err = repo.FindByToAccountID(transaction.AccountFromID, transaction.AccountToID)
		is.Nil(result)
		is.NotNil(err)
		is.Error(err)
	})

	t.Run("should test find all by account to id", func(t *testing.T) {
		repo, mock, transaction := NewTransactionTestMock()
		is := require.New(t)

		row := sqlmock.NewRows([]string{"id", "account_from_id", "amount", "status", "currency", "account_to_id", "created_at", "updated_at", "type", "external_id"}).
			AddRow(transaction.ID, transaction.AccountFromID, transaction.Amount, transaction.Status, transaction.Currency, transaction.AccountToID, transaction.CreatedAt, transaction.UpdatedAt, transaction.Type, transaction.ExternalID)

		page := 1
		limit := 10
		sort := "created_at DESC"

		pagination := &entity.Pagination{
			Page:  page,
			Limit: limit,
			Sort:  sort,
		}

		selectTransaction := fmt.Sprintf(`SELECT * FROM "transactions" WHERE (account_to_id = $1) ORDER BY %s LIMIT %d OFFSET %d`, sort, limit, 0)
		countSelect := `SELECT count(*) FROM "transactions" WHERE (account_to_id = $1)`

		countRow := sqlmock.NewRows([]string{"count"}).AddRow(1)

		mock.ExpectQuery(regexp.QuoteMeta(selectTransaction)).WithArgs(transaction.AccountToID).WillReturnRows(row)
		mock.ExpectQuery(regexp.QuoteMeta(countSelect)).WithArgs(transaction.AccountToID).WillReturnRows(countRow)

		result, total, err := repo.FindAllByToAccountID(transaction.AccountToID, pagination)

		fmt.Println(err)
		is.Nil(err)
		is.Equal(total, 1)
		is.Equal(result[0].ID, transaction.ID)
		is.Equal(result[0].AccountToID, transaction.AccountToID)
		is.Equal(result[0].Amount, transaction.Amount)

		result, total, err = repo.FindAllByToAccountID(transaction.ExternalID, &entity.Pagination{
			Page:  2,
			Limit: limit,
			Sort:  "id DESC",
		})

		is.Nil(result)
		is.Equal(total, 0)
		is.NotNil(err)
	})
}
