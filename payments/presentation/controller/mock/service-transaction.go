package mock

import (
	"github.com/EdlanioJ/kbu/payments/domain/entity"
	"github.com/stretchr/testify/mock"
)

type MockServiceTransactionUseCase struct {
	mock.Mock
}

func NewMockServiceTransactionUseCase() *MockServiceTransactionUseCase {
	return &MockServiceTransactionUseCase{}
}

func (m *MockServiceTransactionUseCase) RegisterServiceTransaction(fromId string, serviceId string, servicePriceId string, amount float64, currency string) (*entity.Transaction, error) {
	args := m.Called(fromId, serviceId, servicePriceId, amount, currency)

	var res0 *entity.Transaction

	if rf, ok := args.Get(0).(func() *entity.Transaction); ok {
		res0 = rf()
	} else {
		if args.Get(0) != nil {
			res0 = args.Get(0).(*entity.Transaction)
		}
	}

	var res1 error
	if rf, ok := args.Get(1).(func() error); ok {
		res1 = rf()
	} else {
		res1 = args.Error(1)
	}
	return res0, res1
}
func (m *MockServiceTransactionUseCase) FindAllByServiceId(serviceId string, page int, limit int, sort string) ([]*entity.Transaction, int, error) {
	args := m.Called(serviceId, page, limit, sort)

	res0 := []*entity.Transaction{}

	if rf, ok := args.Get(0).(func() []*entity.Transaction); ok {
		res0 = rf()
	} else {
		if args.Get(0) != nil {
			res0 = args.Get(0).([]*entity.Transaction)
		}
	}

	var res1 int

	if rf, ok := args.Get(1).(func() int); ok {
		res1 = rf()
	} else {
		res1 = args.Int(1)
	}

	var res2 error
	if rf, ok := args.Get(2).(func() error); ok {
		res2 = rf()
	} else {
		res2 = args.Error(2)
	}

	return res0, res1, res2
}
func (m *MockServiceTransactionUseCase) FindOneByService(serviceId string, transactionId string) (*entity.Transaction, error) {
	args := m.Called(serviceId, transactionId)

	var res0 *entity.Transaction

	if rf, ok := args.Get(0).(func() *entity.Transaction); ok {
		res0 = rf()
	} else {
		if args.Get(0) != nil {
			res0 = args.Get(0).(*entity.Transaction)
		}
	}

	var res1 error
	if rf, ok := args.Get(1).(func() error); ok {
		res1 = rf()
	} else {
		res1 = args.Error(1)
	}
	return res0, res1
}
