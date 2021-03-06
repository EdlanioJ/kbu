package mock

import (
	"github.com/EdlanioJ/kbu/payments/domain/entity"
	"github.com/stretchr/testify/mock"
)

type MockTransactionUseCase struct {
	mock.Mock
}

func NewMockTransactionUseCase() *MockTransactionUseCase {
	return &MockTransactionUseCase{}
}
func (m *MockTransactionUseCase) Register(fromAccount, toAccount, externalID, typeTransaction, currency string, amount float64) (*entity.Transaction, error) {
	args := m.Called(fromAccount, toAccount, externalID, typeTransaction, currency, amount)

	var r0 *entity.Transaction
	if rf, ok := args.Get(0).(func() *entity.Transaction); ok {
		r0 = rf()
	} else {
		if args.Get(0) != nil {
			r0 = args.Get(0).(*entity.Transaction)
		}
	}

	var r1 error
	if rf, ok := args.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = args.Error(1)
	}

	return r0, r1
}

func (m *MockTransactionUseCase) Find(id string) (*entity.Transaction, error) {
	args := m.Called(id)

	var r0 *entity.Transaction
	if rf, ok := args.Get(0).(func() *entity.Transaction); ok {
		r0 = rf()
	} else {
		if args.Get(0) != nil {
			r0 = args.Get(0).(*entity.Transaction)
		}
	}

	var r1 error
	if rf, ok := args.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = args.Error(1)
	}

	return r0, r1
}

func (m *MockTransactionUseCase) FindAll(page int, limit int, sort string) ([]*entity.Transaction, int, error) {
	args := m.Called(page, limit, sort)

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

func (m *MockTransactionUseCase) FindByType(typeTransaction, transactionID string) (*entity.Transaction, error) {
	args := m.Called(typeTransaction, transactionID)
	var r0 *entity.Transaction

	if rf, ok := args.Get(0).(func() *entity.Transaction); ok {
		r0 = rf()
	} else {
		if args.Get(0) != nil {
			r0 = args.Get(0).(*entity.Transaction)
		}
	}

	var r1 error
	if rf, ok := args.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = args.Error(1)
	}

	return r0, r1
}

func (m *MockTransactionUseCase) FindAllByType(typeTransaction string, page int, limit int, sort string) ([]*entity.Transaction, int, error) {
	args := m.Called(typeTransaction, page, limit, sort)

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
func (m *MockTransactionUseCase) FindByExternalID(externalID, transactionID string) (*entity.Transaction, error) {
	args := m.Called(externalID, transactionID)
	var r0 *entity.Transaction

	if rf, ok := args.Get(0).(func() *entity.Transaction); ok {
		r0 = rf()
	} else {
		if args.Get(0) != nil {
			r0 = args.Get(0).(*entity.Transaction)
		}
	}

	var r1 error
	if rf, ok := args.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = args.Error(1)
	}

	return r0, r1
}
func (m *MockTransactionUseCase) FindAllByExternalID(externalID string, page int, limit int, sort string) ([]*entity.Transaction, int, error) {
	args := m.Called(externalID, page, limit, sort)

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
func (m *MockTransactionUseCase) FindAllByFromAccountID(accountID string, page int, limit int, sort string) ([]*entity.Transaction, int, error) {
	args := m.Called(accountID, page, limit, sort)

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

func (m *MockTransactionUseCase) FindByFromAccountID(accountID, transactionID string) (*entity.Transaction, error) {
	args := m.Called(accountID, transactionID)
	var r0 *entity.Transaction

	if rf, ok := args.Get(0).(func() *entity.Transaction); ok {
		r0 = rf()
	} else {
		if args.Get(0) != nil {
			r0 = args.Get(0).(*entity.Transaction)
		}
	}

	var r1 error
	if rf, ok := args.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = args.Error(1)
	}

	return r0, r1
}
func (m *MockTransactionUseCase) FindAllByToAccountID(accountID string, page int, limit int, sort string) ([]*entity.Transaction, int, error) {
	args := m.Called(accountID, page, limit, sort)

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
func (m *MockTransactionUseCase) FindByToAccountID(accountID, transactionID string) (*entity.Transaction, error) {
	args := m.Called(accountID, transactionID)
	var r0 *entity.Transaction

	if rf, ok := args.Get(0).(func() *entity.Transaction); ok {
		r0 = rf()
	} else {
		if args.Get(0) != nil {
			r0 = args.Get(0).(*entity.Transaction)
		}
	}

	var r1 error
	if rf, ok := args.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = args.Error(1)
	}

	return r0, r1
}

func (m *MockTransactionUseCase) Complete(transactionId string) (*entity.Transaction, error) {
	args := m.Called(transactionId)

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

func (m *MockTransactionUseCase) Error(transactionId string) (*entity.Transaction, error) {
	args := m.Called(transactionId)

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
