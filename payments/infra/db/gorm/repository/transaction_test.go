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
	accountFrom, _ := entity.NewAccount(3000)
	accountTo, _ := entity.NewAccount(200)

	transaction, _ := entity.NewTransaction(accountFrom, accountTo, nil, nil, 200, "AKZ")

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

		const insertSql = `INSERT INTO "transactions" ("id","created_at","updated_at","amount","status","currency","account_from_id","account_to_id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "transactions"."id"`
		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(insertSql)).
			WithArgs(
				transaction.ID, transaction.CreatedAt, sqlmock.AnyArg(), transaction.Amount, transaction.Status, transaction.Currency, transaction.AccountFromID, transaction.AccountToID).
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

		const updateSql = `UPDATE "transactions" SET "created_at" = $1, "updated_at" = $2, "amount" = $3, "status" = $4, "currency" = $5, "account_from_id" = $6, "account_to_id" = $7 WHERE "transactions"."id" = $8`
		const selectTransaction = `SELECT * FROM "transactions"  WHERE "transactions"."id" = $1 ORDER BY "transactions"."id" ASC LIMIT 1`

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(updateSql)).
			WithArgs(transaction.CreatedAt, sqlmock.AnyArg(), transaction.Amount, transaction.Status, transaction.Currency, transaction.AccountFromID, transaction.AccountToID, transaction.ID).
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

	t.Run("should test find one by account", func(t *testing.T) {
		repo, mock, transaction := NewTransactionTestMock()
		is := require.New(t)

		row := sqlmock.NewRows([]string{"id", "account_from_id", "amount", "status", "currency", "account_to_id", "created_at", "updated_at"}).
			AddRow(transaction.ID, transaction.AccountFromID, transaction.Amount, transaction.Status, transaction.Currency, transaction.AccountToID, transaction.CreatedAt, transaction.UpdatedAt)

		const selectTransaction = `SELECT * FROM "transactions" WHERE (id = $1 AND account_to_id = $2) ORDER BY "transactions"."id" ASC LIMIT 1`

		mock.ExpectQuery(regexp.QuoteMeta(selectTransaction)).
			WithArgs(transaction.ID, transaction.AccountToID).
			WillReturnRows(row)

		result, err := repo.FindOneByAccount(transaction.ID, transaction.AccountToID)

		is.Nil(err)
		is.Equal(result.ID, transaction.ID)
		is.Equal(result.AccountFromID, transaction.AccountFromID)
		is.Equal(result.Amount, transaction.Amount)

		result, err = repo.FindOneByAccount(uuid.NewV4().String(), transaction.AccountToID)

		is.NotNil(err)
		is.Nil(result)
	})

	t.Run("should test find one by service", func(t *testing.T) {
		repo, mock, _ := NewTransactionTestMock()
		is := require.New(t)

		service, _ := entity.NewService("service", "service description", uuid.NewV4().String(), uuid.NewV4().String())
		accountFrom, _ := entity.NewAccount(3900)
		transaction, _ := entity.NewTransaction(accountFrom, nil, service, nil, 100, "AKZ")

		row := sqlmock.NewRows([]string{"id", "account_from_id", "amount", "status", "currency", "account_to_id", "created_at", "updated_at"}).
			AddRow(transaction.ID, transaction.AccountFromID, transaction.Amount, transaction.Status, transaction.Currency, transaction.AccountToID, transaction.CreatedAt, transaction.UpdatedAt)

		const selectTransaction = `SELECT * FROM "transactions" WHERE (id = $1 AND service_id = $2) ORDER BY "transactions"."id" ASC LIMIT 1`

		mock.ExpectQuery(regexp.QuoteMeta(selectTransaction)).
			WithArgs(transaction.ID, transaction.ServiceID).
			WillReturnRows(row)

		result, err := repo.FindOneByService(transaction.ID, transaction.ServiceID)

		is.Nil(err)
		is.Equal(result.ID, transaction.ID)
		is.Equal(result.AccountFromID, transaction.AccountFromID)
		is.Equal(result.Amount, transaction.Amount)

		result, err = repo.FindOneByService(uuid.NewV4().String(), transaction.ServiceID)

		is.NotNil(err)
		is.Nil(result)
	})

	t.Run("should test find one by store", func(t *testing.T) {
		repo, mock, _ := NewTransactionTestMock()
		is := require.New(t)

		store, _ := entity.NewStore("store", "store description", uuid.NewV4().String(), uuid.NewV4().String())
		accountFrom, _ := entity.NewAccount(3900)
		transaction, _ := entity.NewTransaction(accountFrom, nil, nil, store, 100, "AKZ")

		row := sqlmock.NewRows([]string{"id", "account_from_id", "amount", "status", "currency", "store_id", "created_at", "updated_at"}).
			AddRow(transaction.ID, transaction.AccountFromID, transaction.Amount, transaction.Status, transaction.Currency, transaction.StoreID, transaction.CreatedAt, transaction.UpdatedAt)

		const selectTransaction = `SELECT * FROM "transactions" WHERE (id = $1 AND store_id = $2) ORDER BY "transactions"."id" ASC LIMIT 1`

		mock.ExpectQuery(regexp.QuoteMeta(selectTransaction)).
			WithArgs(transaction.ID, transaction.StoreID).
			WillReturnRows(row)

		result, err := repo.FindOneByStore(transaction.ID, transaction.StoreID)

		is.Nil(err)
		is.Equal(result.ID, transaction.ID)
		is.Equal(result.AccountFromID, transaction.AccountFromID)
		is.Equal(result.Amount, transaction.Amount)

		result, err = repo.FindOneByStore(uuid.NewV4().String(), transaction.StoreID)

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

	t.Run("should test find by account from id", func(t *testing.T) {
		repo, mock, transaction := NewTransactionTestMock()
		is := require.New(t)

		row := sqlmock.NewRows([]string{"id", "account_from_id", "amount", "status", "currency", "account_to_id", "created_at", "updated_at"}).
			AddRow(transaction.ID, transaction.AccountFromID, transaction.Amount, transaction.Status, transaction.Currency, transaction.AccountToID, transaction.CreatedAt, transaction.UpdatedAt)

		page := 1
		limit := 10
		sort := "created_at DESC"

		pagination := &entity.Pagination{
			Page:  page,
			Limit: limit,
			Sort:  sort,
		}

		selectTransaction := fmt.Sprintf(`SELECT * FROM "transactions" WHERE (account_from_id = $1) ORDER BY %s LIMIT %d OFFSET %d`, sort, limit, 0)
		const countSelect = `SELECT count(*) FROM "transactions" WHERE (account_from_id = $1)`

		countRow := sqlmock.NewRows([]string{"count"}).AddRow(1)

		mock.ExpectQuery(regexp.QuoteMeta(selectTransaction)).WithArgs(transaction.AccountFromID).WillReturnRows(row)
		mock.ExpectQuery(regexp.QuoteMeta(countSelect)).WithArgs(transaction.AccountFromID).WillReturnRows(countRow)

		result, total, err := repo.FindByAccountFromId(transaction.AccountFromID, pagination)

		is.Nil(err)
		is.Equal(total, 1)
		is.Equal(result[0].ID, transaction.ID)
		is.Equal(result[0].AccountFromID, transaction.AccountFromID)
		is.Equal(result[0].Amount, transaction.Amount)

		id := uuid.NewV4().String()
		result, total, err = repo.FindByAccountFromId(id, &entity.Pagination{
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

		row := sqlmock.NewRows([]string{"id", "account_from_id", "amount", "status", "currency", "account_to_id", "created_at", "updated_at"}).
			AddRow(transaction.ID, transaction.AccountFromID, transaction.Amount, transaction.Status, transaction.Currency, transaction.AccountToID, transaction.CreatedAt, transaction.UpdatedAt)

		page := 1
		limit := 10
		sort := "created_at DESC"

		pagination := &entity.Pagination{
			Page:  page,
			Limit: limit,
			Sort:  sort,
		}

		selectTransaction := fmt.Sprintf(`SELECT * FROM "transactions" WHERE (account_to_id = $1) ORDER BY %s LIMIT %d OFFSET %d`, sort, limit, page-1)
		countRow := sqlmock.NewRows([]string{"count"}).AddRow(1)
		const countSelect = `SELECT count(*) FROM "transactions" WHERE (account_to_id = $1)`

		mock.ExpectQuery(regexp.QuoteMeta(selectTransaction)).WithArgs(transaction.AccountToID).WillReturnRows(row)
		mock.ExpectQuery(regexp.QuoteMeta(countSelect)).WithArgs(transaction.AccountToID).WillReturnRows(countRow)

		result, total, err := repo.FindByAccountToId(transaction.AccountToID, pagination)

		is.Nil(err)
		is.Equal(total, 1)
		is.Equal(result[0].ID, transaction.ID)
		is.Equal(result[0].AccountFromID, transaction.AccountFromID)
		is.Equal(result[0].Amount, transaction.Amount)

		id := uuid.NewV4().String()
		result, total, err = repo.FindByAccountToId(id, &entity.Pagination{
			Page:  2,
			Limit: 20,
			Sort:  "updated_at DESC",
		})

		is.Equal(total, 0)
		is.Nil(result)
		is.NotNil(err)
	})

	t.Run("should test find by service id", func(t *testing.T) {
		repo, mock, _ := NewTransactionTestMock()
		is := require.New(t)

		service, _ := entity.NewService("service", "service description", uuid.NewV4().String(), uuid.NewV4().String())
		accountFrom, _ := entity.NewAccount(3900)
		transaction, _ := entity.NewTransaction(accountFrom, nil, service, nil, 100, "AKZ")

		page := 1
		limit := 10
		sort := "created_at DESC"

		pagination := &entity.Pagination{
			Page:  page,
			Limit: limit,
			Sort:  sort,
		}

		row := sqlmock.NewRows([]string{"id", "account_from_id", "amount", "status", "currency", "service_id", "created_at", "updated_at"}).
			AddRow(transaction.ID, transaction.AccountFromID, transaction.Amount, transaction.Status, transaction.Currency, transaction.ServiceID, transaction.CreatedAt, transaction.UpdatedAt)

		selectTransaction := fmt.Sprintf(`SELECT * FROM "transactions" WHERE (service_id = $1) ORDER BY %s LIMIT %d OFFSET %d`, sort, limit, page-1)
		const countSelect = `SELECT count(*) FROM "transactions" WHERE (service_id = $1)`
		countRow := sqlmock.NewRows([]string{"count"}).AddRow(1)

		mock.ExpectQuery(regexp.QuoteMeta(selectTransaction)).WithArgs(transaction.ServiceID).WillReturnRows(row)
		mock.ExpectQuery(regexp.QuoteMeta(countSelect)).WithArgs(transaction.ServiceID).WillReturnRows(countRow)

		result, total, err := repo.FindByServiceId(transaction.ServiceID, pagination)

		is.Nil(err)
		is.Equal(total, 1)
		is.Equal(result[0].ID, transaction.ID)
		is.Equal(result[0].AccountFromID, transaction.AccountFromID)
		is.Equal(result[0].Amount, transaction.Amount)

		id := uuid.NewV4().String()
		result, total, err = repo.FindByServiceId(id, &entity.Pagination{
			Page:  10,
			Limit: 2,
			Sort:  "id DESC",
		})

		is.Nil(result)
		is.Equal(total, 0)
		is.NotNil(err)
	})

	t.Run("should test find by store id", func(t *testing.T) {
		repo, mock, _ := NewTransactionTestMock()
		is := require.New(t)

		store, _ := entity.NewStore("store", "store description", uuid.NewV4().String(), uuid.NewV4().String())
		accountFrom, _ := entity.NewAccount(3900)
		transaction, _ := entity.NewTransaction(accountFrom, nil, nil, store, 100, "AKZ")

		page := 1
		limit := 10
		sort := "created_at DESC"

		pagination := &entity.Pagination{
			Page:  page,
			Limit: limit,
			Sort:  sort,
		}

		row := sqlmock.NewRows([]string{"id", "account_from_id", "amount", "status", "currency", "store_id", "created_at", "updated_at"}).
			AddRow(transaction.ID, transaction.AccountFromID, transaction.Amount, transaction.Status, transaction.Currency, transaction.StoreID, transaction.CreatedAt, transaction.UpdatedAt)

		selectTransaction := fmt.Sprintf(`SELECT * FROM "transactions" WHERE (store_id = $1) ORDER BY %s LIMIT %d OFFSET %d`, sort, limit, page-1)
		countRow := sqlmock.NewRows([]string{"count"}).AddRow(1)
		const countSelect = `SELECT count(*) FROM "transactions" WHERE (store_id = $1)`

		mock.ExpectQuery(regexp.QuoteMeta(selectTransaction)).WithArgs(transaction.StoreID).WillReturnRows(row)
		mock.ExpectQuery(regexp.QuoteMeta(countSelect)).WithArgs(transaction.StoreID).WillReturnRows(countRow)

		result, total, err := repo.FindByStoreId(transaction.StoreID, pagination)

		is.Nil(err)
		is.Equal(total, 1)
		is.Equal(result[0].ID, transaction.ID)
		is.Equal(result[0].AccountFromID, transaction.AccountFromID)
		is.Equal(result[0].Amount, transaction.Amount)

		id := uuid.NewV4().String()
		result, total, err = repo.FindByStoreId(id, &entity.Pagination{
			Page:  2,
			Limit: 2,
			Sort:  "id DESC",
		})

		is.Nil(result)
		is.Equal(total, 0)
		is.NotNil(err)
	})
}
