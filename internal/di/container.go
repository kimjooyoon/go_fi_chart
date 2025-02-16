package di

import (
	"fmt"
	"reflect"
	"sync"
)

// Container 의존성 주입 컨테이너 인터페이스입니다.
type Container interface {
	// Register 의존성을 등록합니다.
	Register(name string, constructor interface{}) error
	// Resolve 의존성을 해결합니다.
	Resolve(name string) (interface{}, error)
}

// SimpleContainer 기본 의존성 주입 컨테이너 구현체입니다.
type SimpleContainer struct {
	mu           sync.RWMutex
	constructors map[string]interface{}
	instances    map[string]interface{}
}

// NewSimpleContainer 새로운 SimpleContainer를 생성합니다.
func NewSimpleContainer() *SimpleContainer {
	return &SimpleContainer{
		constructors: make(map[string]interface{}),
		instances:    make(map[string]interface{}),
	}
}

// Register 의존성 생성자를 등록합니다.
func (c *SimpleContainer) Register(name string, constructor interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if constructor == nil {
		return fmt.Errorf("constructor cannot be nil")
	}

	if _, exists := c.constructors[name]; exists {
		return fmt.Errorf("dependency %s already registered", name)
	}

	c.constructors[name] = constructor
	return nil
}

// Resolve 의존성을 해결하고 인스턴스를 반환합니다.
func (c *SimpleContainer) Resolve(name string) (interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 이미 생성된 인스턴스가 있는지 확인
	if instance, exists := c.instances[name]; exists {
		return instance, nil
	}

	// 생성자 찾기
	constructor, exists := c.constructors[name]
	if !exists {
		return nil, fmt.Errorf("dependency %s not registered", name)
	}

	// 생성자 호출
	constructorValue := reflect.ValueOf(constructor)
	if constructorValue.Kind() != reflect.Func {
		return nil, fmt.Errorf("constructor for %s must be a function", name)
	}

	// 생성자 파라미터 해결
	params := make([]reflect.Value, constructorValue.Type().NumIn())
	for i := 0; i < constructorValue.Type().NumIn(); i++ {
		paramType := constructorValue.Type().In(i)
		// TODO: 파라미터 의존성 해결 로직 구현
		params[i] = reflect.New(paramType).Elem()
	}

	// 인스턴스 생성
	results := constructorValue.Call(params)
	if len(results) == 0 {
		return nil, fmt.Errorf("constructor for %s must return at least one value", name)
	}

	instance := results[0].Interface()
	c.instances[name] = instance

	return instance, nil
}
