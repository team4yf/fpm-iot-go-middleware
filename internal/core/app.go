//Package core  核心的结构体
// 用于响应所有的请求
// 处理 设备平台推送的数据，进行整合之后，将消息转发给 MQTT
package core

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/team4yf/fpm-iot-go-middleware/config"
	"github.com/team4yf/fpm-iot-go-middleware/errno"

	s "github.com/team4yf/fpm-iot-go-middleware/internal/service"
	"github.com/team4yf/fpm-iot-go-middleware/pkg/log"
	"github.com/team4yf/fpm-iot-go-middleware/pkg/pubsub"
	"github.com/team4yf/fpm-iot-go-middleware/pkg/utils"
	"github.com/team4yf/fpm-iot-go-middleware/router/middleware"
)

//App the core application
type App struct {
	Router     *mux.Router
	Middleware *middleware.Middleware
	mq         pubsub.PubSub
	Service    s.DeviceService
	m          alice.Chain
}

//Init init the internal component
func (app *App) Init() {

	app.Router = mux.NewRouter()
	app.mq = pubsub.NewMQTTPubSub(config.MqttConfig)
	app.Service = s.NewSimpleDeviceService()
	app.Middleware = &middleware.Middleware{}
	app.m = alice.New(app.Middleware.LoggerMiddleware, app.Middleware.RecoverMiddleware)

}

//ParseBody parse body to struct
func (app *App) ParseBody(r *http.Request, data interface{}) (err error) {
	err = utils.GetBodyStruct(r.Body, &data)

	if err != nil {
		return
	}
	return
}

//Run startup the app
func (app *App) Run(addr string) {

	log.Infof("startup %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, app.Router))
}

//Post add a post handler
func (app *App) Post(url string, handler func(w http.ResponseWriter, r *http.Request)) {
	app.Router.Handle(url, app.m.ThenFunc(handler)).Methods("POST")
}

//Get add a GET handler
func (app *App) Get(url string, handler func(w http.ResponseWriter, r *http.Request)) {
	app.Router.Handle(url, app.m.ThenFunc(handler)).Methods("GET")
}

//Publish publish a message of the topic
func (app *App) Publish(topic string, message []byte) {
	app.mq.Publish(topic, message)
}

//Subscribe subscribe the topic
func (app *App) Subscribe(topic string, handler func(topic, payload interface{})) {
	app.mq.Subscribe(topic, handler)
}

//WriteJSON output json
func (app *App) WriteJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

//SendOk return http:200 and data
func (app *App) SendOk(w http.ResponseWriter, payload interface{}) {
	data, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//Fail return error message but the http:200
func (app *App) Fail(w http.ResponseWriter, result string) {
	err := errno.News(result)
	app.SendError(w, err)
}

//FailWithError return error but the http:200
func (app *App) FailWithError(w http.ResponseWriter, err error) {
	e := errno.NewsWithError(err)
	log.Errorf("error:%v", err)
	app.SendError(w, e)
}

//FailWithCode return error but the http:200
func (app *App) FailWithCode(w http.ResponseWriter, code int, result string) {
	err := errno.NewsWithCode(code, result)
	app.SendError(w, err)
}

//SendError return error but the http:200
func (app *App) SendError(w http.ResponseWriter, err *errno.Errno) {
	data, _ := json.Marshal(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
