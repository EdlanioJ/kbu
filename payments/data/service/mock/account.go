package mock

import (
	"github.com/EdlanioJ/kbu/payments/domain/entity"
	"github.com/stretchr/testify/mock"
)

type MockAccountRepository struct {
	mock.Mock
}

func NewMockAccountRepository() *MockAccountRepository {
	return &MockAccountRepository{}
}

func (mock *MockAccountRepository) Find(id string) (*entity.Account, error) {
	args := mock.Called(id)

	var res0 *entity.Account

	if rf, ok := args.Get(0).(func() *entity.Account); ok {
		res0 = rf()
	} else {
		if args.Get(0) != nil {
			res0 = args.Get(0).(*entity.Account)
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
func (mock *MockAccountRepository) Save(account *entity.Account) error {
	args := mock.Called(account)

	var res0 error
	if rf, ok := args.Get(0).(func() error); ok {
		res0 = rf()
	} else {
		res0 = args.Error(0)
	}

	return res0
}
