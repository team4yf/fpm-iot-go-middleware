package env

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	strData := "FEDC01171CB740E30F000001DE03001C0000010F000001CB000000340000004B000000870000007A0000001400"
	buf := stringToHex(strData)
	device := &JingxunEnvDevice{}
	data, err := device.Parse(buf)
	assert.Nil(t, err, "should not be error")
	fmt.Printf("data: %v", data)
}
