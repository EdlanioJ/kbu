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

func NewServiceTestMock() (*repository.ServiceRepositoryGORM, sqlmock.Sqlmock, *entity.Service) {
	service, _ := entity.NewService("service", "service description", uuid.NewV4().String(), uuid.NewV4().String())

	db, mock, err := sqlmock.New()

	if err != nil {
		panic(err)
	}

	gdb, err := gorm.Open("postgres", db)

	gdb.LogMode(false)
	if err != nil {
		panic(err)
	}

	repo := repository.NewServiceRepository(gdb)
	return repo, mock, service
}

func TestServiceRepository(t *testing.T) {
	t.Parallel()

	t.Run("should test find", func(t *testing.T) {
		repo, mock, service := NewServiceTestMock()
		is := require.New(t)

		row := sqlmock.NewRows([]string{"id", "description", "name", "from_id", "type_id", "created_at", "updated_at"}).
			AddRow(service.ID, service.Description, service.Name, service.FromID, service.TypeID, service.CreatedAt, service.UpdatedAt)

		const sqlService = `SELECT * FROM "services" WHERE (id = $1) ORDER BY "services"."id" ASC LIMIT 1`

		mock.ExpectQuery(regexp.QuoteMeta(sqlService)).
			WithArgs(service.ID).
			WillReturnRows(row)

		result, err := repo.Find(service.ID)

		is.Nil(err)
		is.Equal(result.ID, service.ID)
		is.Equal(service.Description, result.Description)

		result, err = repo.Find(uuid.NewV4().String())

		is.NotNil(err)
		is.Nil(result)
	})
	t.Run("should test find service by id and status", func(t *testing.T) {
		repo, mock, service := NewServiceTestMock()
		is := require.New(t)

		row := sqlmock.NewRows([]string{"id", "description", "name", "from_id", "type_id", "created_at", "updated_at"}).
			AddRow(service.ID, service.Description, service.Name, service.FromID, service.TypeID, service.CreatedAt, service.UpdatedAt)

		const sqlService = `SELECT * FROM "services" WHERE (id = $1 AND status = $2) ORDER BY "services"."id" ASC LIMIT 1`

		mock.ExpectQuery(regexp.QuoteMeta(sqlService)).
			WithArgs(service.ID, service.Status).
			WillReturnRows(row)

		result, err := repo.FindServiceByIdAndStatus(service.ID, service.Status)

		is.Nil(err)
		is.Equal(result.ID, service.ID)
		is.Equal(service.Description, result.Description)

		result, err = repo.FindServiceByIdAndStatus(uuid.NewV4().String(), service.Status)

		is.NotNil(err)
		is.Nil(result)
	})
}
