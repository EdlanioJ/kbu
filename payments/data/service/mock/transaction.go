package mock

import (
	"github.com/EdlanioJ/kbu/payments/domain/entity"
	"github.com/stretchr/testify/mock"
)

type MockTransactionRepository struct {
	mock.Mock
}

func NewMockTransactionRepository() *MockTransactionRepository {
	return &MockTransactionRepository{}
}

func (m *MockTransactionRepository) Register(transaction *entity.Transaction) error {

	args := m.Called(transaction)

	var res0 error
	if rf, ok := args.Get(0).(func() error); ok {
		res0 = rf()
	} else {
		res0 = args.Error(0)
	}

	return res0
}

func (m *MockTransactionRepository) Save(transaction *entity.Transaction) error {
	args := m.Called(transaction)

	var res0 error
	if rf, ok := args.Get(0).(func() error); ok {
		res0 = rf()
	} else {
		res0 = args.Error(0)
	}

	return res0
}

func (m *MockTransactionRepository) Find(id string) (*entity.Transaction, error) {
	args := m.Called(id)

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

func (m *MockTransactionRepository) FindAll(pagination *entity.Pagination) ([]*entity.Transaction, int, error) {
	args := m.Called(pagination)

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

func (m *MockTransactionRepository) FindByType(transactionID, transactionType string) (*entity.Transaction, error) {
	args := m.Called(transactionID, transactionType)

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

func (m *MockTransactionRepository) FindAllByType(transactionType string, pagination *entity.Pagination) ([]*entity.Transaction, int, error) {
	args := m.Called(transactionType, pagination)

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

func (m *MockTransactionRepository) FindByExternalID(transactionID, externalID string) (*entity.Transaction, error) {
	args := m.Called(transactionID, externalID)

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

func (m *MockTransactionRepository) FindAllByExternalID(externalID string, pagination *entity.Pagination) ([]*entity.Transaction, int, error) {
	args := m.Called(externalID, pagination)
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

func (m *MockTransactionRepository) FindByFromAccountID(transactionID, accountID string) (*entity.Transaction, error) {
	args := m.Called(transactionID, accountID)

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

func (m *MockTransactionRepository) FindAllByFromAccountID(accountID string, pagination *entity.Pagination) ([]*entity.Transaction, int, error) {
	args := m.Called(accountID, pagination)

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

func (m *MockTransactionRepository) FindByToAccountID(transactionID, accountID string) (*entity.Transaction, error) {
	args := m.Called(transactionID, accountID)

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

func (m *MockTransactionRepository) FindAllByToAccountID(accountID string, pagination *entity.Pagination) ([]*entity.Transaction, int, error) {
	args := m.Called(accountID, pagination)
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
