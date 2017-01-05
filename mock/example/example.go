package mock

import (
	"errors"
	"fmt"
)

type Object interface {
	DoSomething(number int) (bool, error)
	DoSomething2(name string) (ReturnObject, error)
}

type MyObject struct {
}

type ReturnObject struct {
	name string
}

func (m *MyObject) DoSomething(number int) (bool, error) {
	if number > 0 {
		return true, nil
	}
	return false, errors.New("game over")
}

func (m *MyObject) DoSomething2(name string) (ReturnObject, error) {
	return ReturnObject{name}, nil
}

func targetFuncThatDoesSomethingWithObj(m Object, number int) {
	result, err := m.DoSomething(number)
	fmt.Println(result, err)
}

func targetFuncThatDoesSomethingWithObj2(m Object, name string) {
	result, err := m.DoSomething2(name)
	fmt.Println(result, err)
}
