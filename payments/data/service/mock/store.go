package mock

import (
	"github.com/EdlanioJ/kbu/payments/domain/entity"
	"github.com/stretchr/testify/mock"
)

type MockStoreRepository struct {
	mock.Mock
}

func NewMockStoreRepository() *MockStoreRepository {
	return &MockStoreRepository{}
}

func (m *MockStoreRepository) Find(id string) (*entity.Store, error) {
	args := m.Called(id)

	var res0 *entity.Store

	if rf, ok := args.Get(0).(func() *entity.Store); ok {
		res0 = rf()
	} else {
		if args.Get(0) != nil {
			res0 = args.Get(0).(*entity.Store)
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

func (m *MockStoreRepository) FindStoreByIdAndStatus(id string, status string) (*entity.Store, error) {
	args := m.Called(id, status)

	result := args.Get(0)

	return result.(*entity.Store), args.Error(1)
}
