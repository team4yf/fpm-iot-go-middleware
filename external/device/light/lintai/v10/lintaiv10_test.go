package lintaiv10

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientInit(t *testing.T) {
	options := &Options{
		AppID:       "LT0314fbf27a4d2986",
		AppSecret:   "1bc7b874c74623298a6",
		Username:    "18796664408",
		TokenExpire: 60 * 1000 * 24 * 7,

		Enviroment: "prod",
		BaseURL:    "http://101.132.142.5:8088/api",
	}

	client := NewClient(options)

	err := client.Init()

	assert.Nil(t, err, "should not be err")
}
