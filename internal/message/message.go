//Package message the message data struct defination
package message

//Header the header of the message
type Header struct {
	Version   int    `json:"v"`
	NameSpace string `json:"ns"`
	Name      string `json:"name"`
	AppID     string `json:"appId"`
	ProjID    int64  `json:"projId"`
	Source    string `json:"source"`
}

//Device the device data struct of the payload
type Device struct {
	ID      string                 `json:"id"`
	Type    string                 `json:"type"`
	Name    string                 `json:"name"`
	Brand   string                 `json:"brand"`
	Version string                 `json:"v"`
	Extra   map[string]interface{} `json:"x,omitempty"`
}

//D2SPayload the payload data struct
type D2SPayload struct {
	Device    *Device     `json:"device"`
	Data      interface{} `json:"data"`
	Cgi       string      `json:"cgi"`
	Timestamp int64       `json:"timestamp"`
}

//D2SFeedback the body of the feedback message
type D2SFeedback struct {
	Result    interface{} `json:"result"`
	MsgID     string      `json:"msgId"`
	Cgi       string      `json:"cgi"`
	Timestamp int64       `json:"timestamp"`
}

//S2DPayload the payload data struct
type S2DPayload struct {
	Device    *Device     `json:"device"`
	Argument  interface{} `json:"arg"`
	MsgID     string      `json:"msgId"`
	NetID     string      `json:"netId"`
	Cmd       string      `json:"cmd"`
	Cgi       string      `json:"cgi"`
	Feedback  int         `json:"feedback,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

//D2SMessage device to server message
type D2SMessage struct {
	Header  *Header     `json:"header"`
	Payload *D2SPayload `json:"payload"`
}

//S2DMessage server to device message
type S2DMessage struct {
	Header  *Header                `json:"header"`
	Bind    map[string]interface{} `json:"bind,omitempty"`
	Payload []*S2DPayload          `json:"payload"`
}

//D2SFeedbackMessage feedback message
type D2SFeedbackMessage struct {
	Header   *Header      `json:"header"`
	Feedback *D2SFeedback `json:"feedback"`
}

//EnvPayload env data struct
// "min_wind_dir": 0,
// "avg_wind_dir": 0,
// "max_wind_dir": 0,
// "min_wind_speed": 0,
// "avg_wind_speed": 0,
// "max_wind_speed": 2,
// "temp":253,//温度值 （-550~1250）精度0.1
// "humidity": 597,
// "pressure": 10087,
// "rainfall": 0,
// "radiation": 0,
// "u_rays": 0,
// "noise": 562,
// "pm2_5": 0,
// "pm10": 0
type EnvPayload struct {
	Temp         int `json:"temp"`
	MinWindDir   int `json:"min_wind_dir"`
	AvgWindDir   int `json:"avg_wind_dir"`
	MaxWindDir   int `json:"max_wind_dir"`
	MinWindSpeed int `json:"min_wind_speed"`
	AvgWindSpeed int `json:"avg_wind_speed"`
	MaxWindSpeed int `json:"max_wind_speed"`
	Humidity     int `json:"humidity"`
	Pressure     int `json:"pressure"`
	Radiation    int `json:"radiation"`
	Rainfall     int `json:"rainfall"`
	URays        int `json:"u_rays"`
	Noise        int `json:"noise"`
	PM25         int `json:"pm2_5"`
	PM10         int `json:"pm10"`
}
