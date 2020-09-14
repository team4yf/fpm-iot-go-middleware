//Package env about the env
package env

import (
	"encoding/hex"
	"fmt"

	"github.com/team4yf/fpm-iot-go-middleware/internal/message"
)

//EnvDevice the api around the env
type EnvDevice interface {
	Version() string
	Parse(data *[]byte) (*message.D2SMessage, error)
	Init() error
}

//NewEnvDevice create new env device accroding by the brand
func NewEnvDevice(brand string) (device EnvDevice) {
	switch brand {
	case "jingxun":
		device = &JingxunEnvDevice{}
	}
	if device == nil {
		panic(fmt.Errorf("unknown brand: %s", brand))
	}
	device.Init()
	return device
}

func stringToHex(str string) *[]byte {
	data, err := hex.DecodeString(str)
	if err != nil {
		// handle error
		return new([]byte)
	}
	return &data
}
