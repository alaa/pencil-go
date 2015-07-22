package registry

import (
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestSynchronizeWhenNoServicesWereRegisteredBefore(t *testing.T) {
	serviceRepository := new(MockServiceRepository)
	containerRepository := new(MockContainerRepository)
	registry := NewRegistry(containerRepository, serviceRepository)

	serviceRepository.On("GetAllIds").Return([]string{})
	serviceRepository.On("Register", &Service{
		ID:      "bd1d34c0ebeeb62dfdcc57327aca15d2ef3cbc39a60e44aecb7085a8d1f89fd9",
		Service: "/elated_kirch",
		Port:    22,
	}).Return(nil)
	serviceRepository.On("Register", &Service{
		ID:      "f717f795bcccd674628b92f77a72f4b80b2c6b5da289846a0edbd21fb4c462db",
		Service: "/naughty_heisenberg",
		Port:    9000,
	}).Return(nil)

	containerRepository.On("GetAll").Return(
		[]*Container{
			&Container{
				ID:   "bd1d34c0ebeeb62dfdcc57327aca15d2ef3cbc39a60e44aecb7085a8d1f89fd9",
				Name: "/elated_kirch",
				Port: 22,
			},
			&Container{
				ID:   "f717f795bcccd674628b92f77a72f4b80b2c6b5da289846a0edbd21fb4c462db",
				Name: "/naughty_heisenberg",
				Port: 9000,
			},
		},
	)

	registry.Synchronize()

	serviceRepository.AssertExpectations(t)
	containerRepository.AssertExpectations(t)
}

func TestSynchronieWhenAllServicesWereRegisteredBefore(t *testing.T) {
	serviceRepository := new(MockServiceRepository)
	containerRepository := new(MockContainerRepository)
	registry := NewRegistry(containerRepository, serviceRepository)

	serviceRepository.On("GetAllIds").Return([]string{
		"bd1d34c0ebeeb62dfdcc57327aca15d2ef3cbc39a60e44aecb7085a8d1f89fd9",
		"f717f795bcccd674628b92f77a72f4b80b2c6b5da289846a0edbd21fb4c462db",
	})
	containerRepository.AssertNotCalled(t, "Register")

	containerRepository.On("GetAll").Return(
		[]*Container{
			&Container{
				ID:   "bd1d34c0ebeeb62dfdcc57327aca15d2ef3cbc39a60e44aecb7085a8d1f89fd9",
				Name: "/elated_kirch",
				Port: 22,
			},
			&Container{
				ID:   "f717f795bcccd674628b92f77a72f4b80b2c6b5da289846a0edbd21fb4c462db",
				Name: "/naughty_heisenberg",
				Port: 9000,
			},
		},
	)

	registry.Synchronize()

	serviceRepository.AssertExpectations(t)
	containerRepository.AssertExpectations(t)
}

func TestSynchronieWhenOneServiceIsMissingAndOneIsRedundant(t *testing.T) {
	serviceRepository := new(MockServiceRepository)
	containerRepository := new(MockContainerRepository)
	registry := NewRegistry(containerRepository, serviceRepository)

	serviceRepository.On("GetAllIds").Return([]string{
		"bd1d34c0ebeeb62dfdcc57327aca15d2ef3cbc39a60e44aecb7085a8d1f89fd9",
		"0g1d34c0ebeeb62dfdcc57327aca15d2ef3cbc39a60e44aecb7085a8d1f89fd9",
	})
	containerRepository.AssertNotCalled(t, "Register")

	containerRepository.On("GetAll").Return(
		[]*Container{
			&Container{
				ID:   "bd1d34c0ebeeb62dfdcc57327aca15d2ef3cbc39a60e44aecb7085a8d1f89fd9",
				Name: "/elated_kirch",
				Port: 22,
			},
			&Container{
				ID:   "f717f795bcccd674628b92f77a72f4b80b2c6b5da289846a0edbd21fb4c462db",
				Name: "/naughty_heisenberg",
				Port: 9000,
			},
		},
	)

	serviceRepository.On("Register", &Service{
		ID:      "f717f795bcccd674628b92f77a72f4b80b2c6b5da289846a0edbd21fb4c462db",
		Service: "/naughty_heisenberg",
		Port:    9000,
	}).Return(nil)
	serviceRepository.On("Unregister", "0g1d34c0ebeeb62dfdcc57327aca15d2ef3cbc39a60e44aecb7085a8d1f89fd9").Return(nil)

	registry.Synchronize()

	serviceRepository.AssertExpectations(t)
	containerRepository.AssertExpectations(t)
}

type MockServiceRepository struct {
	mock.Mock
}

type MockContainerRepository struct {
	mock.Mock
}

func (msr *MockServiceRepository) GetAllIds() []string {
	args := msr.Called()
	return args.Get(0).([]string)

}
func (msr *MockServiceRepository) Register(service *Service) error {
	args := msr.Called(service)
	return args.Error(0)
}

func (msr *MockServiceRepository) Unregister(serviceID string) error {
	args := msr.Called(serviceID)
	return args.Error(0)
}

func (mcr *MockContainerRepository) GetAll() []*Container {
	args := mcr.Called()
	return args.Get(0).([]*Container)
}
