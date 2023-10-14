package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	s := NewStack(128)

	s.Push(1)
	s.Push(2)

	value := s.Pop()
	assert.Equal(t, value, 1)

	value = s.Pop()
	assert.Equal(t, value, 2)
}

func TestVM(t *testing.T) {
	data := []byte{0x02, 0x0a, 0x03, 0x0a, 0x0b, 0x4f, 0x0c, 0x4f, 0x0c, 0x46, 0x0c, 0x03, 0x0a, 0x0d, 0x0f}
	pushFoo := []byte{0x4f, 0x0c, 0x4f, 0x0c, 0x46, 0x0c, 0x03, 0x0a, 0x0d, 0xae}
	
	data = append(data, pushFoo...)

	contractState := NewState()
	vm := NewVM(data, contractState)
	assert.Nil(t, vm.Run())

	//fmt.Printf("%+v", vm.stack.data)
	value := vm.stack.Pop().([]byte)
	valueSerialized := deserializeInt64(value)
	
	assert.Equal(t, valueSerialized, int64(5))
	
	// valueBytes, err := contractState.Get([]byte("FOO"))
	// assert.Nil(t, err)
	// assert.Equal(t, value, int64(5))
}

func TestVmMul(t *testing.T) {
	data := []byte{0x02, 0x0a, 0x02, 0x0a, 0xea}
	contractState := NewState()
	vm := NewVM(data, contractState)
	assert.Nil(t, vm.Run())

	result := vm.stack.Pop().(int)
	assert.Equal(t, result, 4)
}

func TestVmDiv(t *testing.T) {
	data := []byte{0x04, 0x0a, 0x02, 0x0a, 0xfd}
	contractState := NewState()
	vm := NewVM(data, contractState)
	assert.Nil(t, vm.Run())

	result := vm.stack.Pop().(int)
	assert.Equal(t, result, 2)
}