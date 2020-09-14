package env

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/team4yf/fpm-iot-go-middleware/internal/message"
	"github.com/team4yf/fpm-iot-go-middleware/internal/service"
)

//JingxunEnvDevice the jingxun env device.
type JingxunEnvDevice struct {
	deviceService service.DeviceService
}

func (device *JingxunEnvDevice) Init() error {
	device.deviceService = service.GetSimpleDeviceService()
	return nil
}

func (device *JingxunEnvDevice) Version() string {
	return "beta"
}

//Parse decode the tcp data
//FEDC01171CB740E30F000001DE03001C0000010F000001CB000000340000004B000000870000007A0000001400
//FEDC 01 171CB740E30F 00000178 03 001C 0000010F 温度 000001CB 湿度 00000034 PM2.5 0000004B PM10 00000087 信号强度 0000007A 错误码 00000014 版本号 00
// 0  2 protocol
// 2  1 function code
// 3  6 device id
// 9  4 00000178 ?
// 13 1 03 ?
// 14 2 001C ?
// 16 4 0000010F 温度
// 20 4 000001CB 湿度
// 24 4 00000034 PM2.5
// 28 4 0000004B PM10
// 32 4 00000087 信号强度
// 36 4 0000007A 错误码
// 40 4 00000014 版本号
func (device *JingxunEnvDevice) Parse(data *[]byte) (payload *message.D2SMessage, err error) {
	if (*data)[0] != 0xfe || (*data)[1] != 0xdc {
		err = fmt.Errorf("data protocol should be [fedc]")
		return
	}
	sid := (*data)[3:9]
	temp := (*data)[16:20]
	hu := (*data)[20:24]
	pm25 := (*data)[24:28]
	pm10 := (*data)[28:32]

	envPayload := &message.EnvPayload{}
	envPayload.Temp = byteToInt32(&temp)
	envPayload.Humidity = byteToInt32(&hu)
	envPayload.PM25 = byteToInt32(&pm25)
	envPayload.PM10 = byteToInt32(&pm10)

	d2sPayload := &message.D2SPayload{}
	d2sPayload.Data = envPayload
	d2sPayload.Device = &message.Device{
		ID:      fmt.Sprintf("%X", sid),
		Type:    "env",
		Brand:   "jingxun",
		Version: device.Version(),
	}
	d2sPayload.Cgi = d2sPayload.Device.ID
	d2sPayload.Timestamp = time.Now().Unix()

	// 通过设备和id获取到具体对应的项目信息，如果设备不存在或者设备状态不对的话，会抛出异常信息
	uuid, projid, err := device.deviceService.Receive("env", "jingxun", "beat", d2sPayload.Cgi)
	if err != nil {
		return
	}

	msgHeader := &message.Header{
		Version:   10,
		NameSpace: "FPM.Lamp.env",
		Name:      "beat",
		AppID:     uuid,
		ProjID:    projid,
		Source:    "TCP",
	}
	payload = &message.D2SMessage{
		Header:  msgHeader,
		Payload: d2sPayload,
	}

	return
}

func byteToInt32(b *[]byte) int {
	bytesBuffer := bytes.NewBuffer(*b)
	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return int(x)
}
