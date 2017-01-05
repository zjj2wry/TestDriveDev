package mock

import "github.com/stretchr/testify/mock"

var _ Object = new(Mock)

type Mock struct {
	mock.Mock
}

func (m *Mock) DoSomething(number int) (bool, error) {
	args := m.Called(number)
	return args.Bool(0), args.Error(1)

}

func (m *Mock) DoSomething2(name string) (ReturnObject, error) {
	args := m.Called(name)
	return args.Get(0).(ReturnObject), args.Error(1)
}
