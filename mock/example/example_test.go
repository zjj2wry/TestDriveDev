package mock

import (
	"fmt"
	"testing"
)

func TestSomething(t *testing.T) {
	mock := new(Mock)
	call := mock.On("DoSomething", 123).Return(true, nil)
	fmt.Println(call.Method, call.Arguments, call.Parent, call.ReturnArguments, call.RunFn)
	targetFuncThatDoesSomethingWithObj(mock, 123)

	mock.AssertExpectations(t)
}

func TestSomething2(t *testing.T) {
	mock := new(Mock)
	mock.On("DoSomething2", "hello").Return(ReturnObject{"world"}, nil)
	targetFuncThatDoesSomethingWithObj2(mock, "hello")

	// mock.AssertExpectations(t)
}
