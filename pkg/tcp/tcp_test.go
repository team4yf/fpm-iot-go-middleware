package tcp

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCreateTCPServer(t *testing.T){
	receiver := NewNetReceiver(func(topic string, data []byte){
		assert.Equal(t, "test", topic, "should be test")
	})

	receiver.Listen(5001)
	

}