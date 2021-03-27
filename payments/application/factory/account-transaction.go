package factory

import (
	"github.com/EdlanioJ/kbu/payments/data/service"
	"github.com/EdlanioJ/kbu/payments/infra/db/gorm/repository"
	"github.com/EdlanioJ/kbu/payments/presentation/controller"
	"github.com/jinzhu/gorm"
)

func AccountTransactionControllerFactory(database *gorm.DB) *controller.AccountTransaction {
	transactionRepository := repository.NewTransactionRepository(database)
	accountRepository := repository.NewAccountRepository(database)

	transactionService := service.NewAccountTransaction(transactionRepository, accountRepository)

	return controller.NewAccountTransaction(transactionService)
}
