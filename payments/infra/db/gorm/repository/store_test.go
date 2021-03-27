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

func NewStoreTestMock() (*repository.StoreRepositoryGORM, sqlmock.Sqlmock, *entity.Store) {
	store, _ := entity.NewStore("store", "store description", uuid.NewV4().String(), uuid.NewV4().String())
	db, mock, err := sqlmock.New()

	if err != nil {
		panic(err)
	}

	gdb, err := gorm.Open("postgres", db)

	gdb.LogMode(false)
	if err != nil {
		panic(err)
	}
	repo := repository.NewStoreRepository(gdb)
	return repo, mock, store
}

func TestStoreRepository(t *testing.T) {
	t.Parallel()

	t.Run("should test find", func(t *testing.T) {
		repo, mock, store := NewStoreTestMock()
		is := require.New(t)

		row := sqlmock.NewRows([]string{"id", "name", "description", "from_id", "type_id", "created_at", "updated_at"}).
			AddRow(store.ID, store.Name, store.Description, store.FromID, store.TypeID, store.CreatedAt, store.UpdatedAt)

		const sqlSelect = `SELECT * FROM "stores" WHERE (id = $1) ORDER BY "stores"."id" ASC LIMIT 1`

		mock.ExpectQuery(regexp.QuoteMeta(sqlSelect)).
			WithArgs(store.ID).
			WillReturnRows(row)

		result, err := repo.Find(store.ID)

		is.Nil(err)
		is.Equal(result.ID, store.ID)

		result, err = repo.Find(uuid.NewV4().String())

		is.NotNil(err)
		is.Nil(result)
	})

	t.Run("should test find store by id and status", func(t *testing.T) {
		repo, mock, store := NewStoreTestMock()
		is := require.New(t)

		row := sqlmock.NewRows([]string{"id", "name", "description", "from_id", "type_id", "created_at", "updated_at"}).
			AddRow(store.ID, store.Name, store.Description, store.FromID, store.TypeID, store.CreatedAt, store.UpdatedAt)

		const sqlSelect = `SELECT * FROM "stores" WHERE (id = $1 AND status = $2) ORDER BY "stores"."id" ASC LIMIT 1`

		mock.ExpectQuery(regexp.QuoteMeta(sqlSelect)).
			WithArgs(store.ID, store.Status).
			WillReturnRows(row)

		result, err := repo.FindStoreByIdAndStatus(store.ID, store.Status)

		is.Nil(err)
		is.Equal(result.ID, store.ID)

		result, err = repo.FindStoreByIdAndStatus(uuid.NewV4().String(), store.Status)

		is.NotNil(err)
		is.Error(err)
		is.Nil(result)
	})
}
