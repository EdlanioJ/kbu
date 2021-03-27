package factory

import (
	"github.com/EdlanioJ/kbu/payments/data/service"
	"github.com/EdlanioJ/kbu/payments/infra/db/gorm/repository"
	"github.com/EdlanioJ/kbu/payments/presentation/controller"
	"github.com/jinzhu/gorm"
)

func StoreTransactionControllerFactory(database *gorm.DB) *controller.StoreTransaction {
	transactionRepository := repository.NewTransactionRepository(database)
	accountRepository := repository.NewAccountRepository(database)
	storeRepository := repository.NewStoreRepository(database)

	serviceTransaction := service.NewStoreTransaction(
		accountRepository,
		storeRepository,
		transactionRepository,
	)

	return controller.NewStoreTransaction(serviceTransaction)
}
