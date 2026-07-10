package tools

import (
	"fmt"
	"sync"
	"turtle/core/lgr"
)

type ConstructFunction func() interface{}

type ClassFactory struct {
	constructors map[string]ConstructFunction
}

var (
	factoryInstance *ClassFactory
	once            sync.Once
)

func Instance() *ClassFactory {
	once.Do(func() {
		factoryInstance = &ClassFactory{
			constructors: make(map[string]ConstructFunction),
		}
	})
	return factoryInstance
}

func (f *ClassFactory) RegisterConstructor(typeName string, constructor ConstructFunction) {

	_, ok := f.constructors[typeName]

	if ok {
		lgr.Error("Entity [", typeName, "] in factory already exits")
	} else {
		f.constructors[typeName] = constructor
	}

}

func (f *ClassFactory) HasContructor(typeName string) bool {
	_, ok := f.constructors[typeName]
	return ok
}

func RegisterClass[T any](typeName string, constructor func() *T) {
	Instance().RegisterConstructor(typeName, func() interface{} {
		return constructor()
	})
}

func Construct[T any](typeName string) (*T, error) {
	constructor, exists := Instance().constructors[typeName]
	if !exists {
		return nil, fmt.Errorf("unable to find constructor for type [%s]", typeName)
	}

	instance, ok := constructor().(*T)
	if !ok {
		return nil, fmt.Errorf("unable to cast constructed instance to type [%s]", typeName)
	}

	return instance, nil
}

func ConstructAny(typeName string) (*any, error) {
	constructor, exists := Instance().constructors[typeName]
	if !exists {
		return nil, fmt.Errorf("unable to find constructor for type [%s]", typeName)
	}

	instance := constructor()
	if instance == nil {

		return nil, fmt.Errorf("unable to cast constructed instance to type [%s]", typeName)
	}

	return &instance, nil
}
