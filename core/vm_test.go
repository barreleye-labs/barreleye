package core

import (
	"fmt"
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
	dataOther := []byte{0x02, 0x0a, 0x03, 0x0a, 0x0b, 0x4d, 0x0c, 0x4f, 0x0c, 0x46, 0x0c, 0x03, 0x0a, 0x0d, 0x0f}
	
	data = append(data, dataOther...)
	
	contractState := NewState()
	vm := NewVM(data, contractState)
	assert.Nil(t, vm.Run())

	fmt.Printf("%+v\n", contractState)

	valueBytes, err := contractState.Get([]byte("FOO"))
	assert.Nil(t, err)
	value := deserializeInt64(valueBytes)
	assert.Equal(t, value, int64(5))
}
