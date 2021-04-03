package factory

import (
	"github.com/EdlanioJ/kbu/payments/data/service"
	"github.com/EdlanioJ/kbu/payments/infra/db/gorm/repository"
	"github.com/EdlanioJ/kbu/payments/presentation/controller"
	"github.com/jinzhu/gorm"
)

func TransactionControllerFactory(database *gorm.DB) *controller.Transaction {
	transactionRepo := repository.NewTransactionRepository(database)
	accountRepo := repository.NewAccountRepository(database)
	transactionService := service.NewTransaction(transactionRepo, accountRepo)

	return controller.NewTransaction(transactionService)
}
