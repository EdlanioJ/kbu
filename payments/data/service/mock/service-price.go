package mock

import (
	"github.com/EdlanioJ/kbu/payments/domain/entity"
	"github.com/stretchr/testify/mock"
)

type MockServicePriceRepository struct {
	mock.Mock
}

func NewMockServicePriceRepository() *MockServicePriceRepository {
	return &MockServicePriceRepository{}
}

func (m *MockServicePriceRepository) Find(id string) (*entity.ServicePrice, error) {
	args := m.Called(id)

	var res0 *entity.ServicePrice
	if rf, ok := args.Get(0).(func() *entity.ServicePrice); ok {
		res0 = rf()
	} else {
		if args.Get(0) != nil {
			res0 = args.Get(0).(*entity.ServicePrice)
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
