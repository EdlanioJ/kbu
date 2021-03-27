package factory

import (
	"github.com/EdlanioJ/kbu/payments/data/service"
	"github.com/EdlanioJ/kbu/payments/infra/db/gorm/repository"
	"github.com/EdlanioJ/kbu/payments/presentation/controller"
	"github.com/jinzhu/gorm"
)

func ServiceTransactionControllerFactory(database *gorm.DB) *controller.ServiceTransaction {
	transactionRepository := repository.NewTransactionRepository(database)
	accountRepository := repository.NewAccountRepository(database)
	serviceRepository := repository.NewServiceRepository(database)
	servicePriceRepository := repository.NewServicePriceRepository(database)

	serviceTransaction := service.NewServiceTransaction(
		accountRepository,
		servicePriceRepository,
		serviceRepository,
		transactionRepository,
	)

	return controller.NewServiceTransaction(serviceTransaction)
}
