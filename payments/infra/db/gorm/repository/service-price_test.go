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

func NewServicePriceTestMock() (*repository.ServicePriceRepositoryGORM, sqlmock.Sqlmock, *entity.ServicePrice) {
	service, _ := entity.NewService("service", "service description", uuid.NewV4().String(), uuid.NewV4().String())
	servicePrice, _ := entity.NewServicePrice(service, "service price description", 10, "AOA")

	db, mock, err := sqlmock.New()

	if err != nil {
		panic(err)
	}

	gdb, err := gorm.Open("postgres", db)
	gdb.LogMode(false)
	if err != nil {
		panic(err)
	}

	repo := repository.NewServicePriceRepository(gdb)

	return repo, mock, servicePrice
}

func TestServicePriceRepository(t *testing.T) {
	t.Parallel()

	t.Run("should test find", func(t *testing.T) {
		repo, mock, servicePrice := NewServicePriceTestMock()
		is := require.New(t)

		row := sqlmock.NewRows([]string{"id", "description", "amount", "service_id", "currency", "created_at", "updated_at"}).
			AddRow(servicePrice.ID, servicePrice.Description, servicePrice.Amount, servicePrice.ServiceID, servicePrice.Currency, servicePrice.CreatedAt, servicePrice.UpdatedAt)

		const sql = `SELECT * FROM "service_prices" WHERE (id = $1) ORDER BY "service_prices"."id" ASC LIMIT 1`

		mock.ExpectQuery(regexp.QuoteMeta(sql)).
			WithArgs(servicePrice.ID).
			WillReturnRows(row)

		result, err := repo.Find(servicePrice.ID)

		is.Nil(err)
		is.Equal(result.ID, servicePrice.ID)
		is.Equal(result.Amount, servicePrice.Amount)
		is.Equal(result.ServiceID, servicePrice.ServiceID)
		is.Equal(result.CreatedAt, servicePrice.CreatedAt)
		is.Equal(result.Currency, servicePrice.Currency)

		result, err = repo.Find(uuid.NewV4().String())

		is.NotNil(err)
		is.Error(err)
		is.Nil(result)
	})
}
