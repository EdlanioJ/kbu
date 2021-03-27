package mock

import (
	"github.com/EdlanioJ/kbu/payments/domain/entity"
	"github.com/stretchr/testify/mock"
)

type MockServiceRepository struct {
	mock.Mock
}

func NewMockServiceRepository() *MockServiceRepository {
	return &MockServiceRepository{}
}

func (m *MockServiceRepository) Find(id string) (*entity.Service, error) {
	args := m.Called(id)

	var res0 *entity.Service

	if rf, ok := args.Get(0).(func() *entity.Service); ok {
		res0 = rf()
	} else {
		if args.Get(0) != nil {
			res0 = args.Get(0).(*entity.Service)
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
func (m *MockServiceRepository) FindServiceByIdAndStatus(id string, status string) (*entity.Service, error) {
	args := m.Called(id, status)

	var res0 *entity.Service

	if rf, ok := args.Get(0).(func() *entity.Service); ok {
		res0 = rf()
	} else {
		if args.Get(0) != nil {
			res0 = args.Get(0).(*entity.Service)
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
