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
	"github.com/team4yf/fpm-iot-go-middleware/external/rest"

	s "github.com/team4yf/fpm-iot-go-middleware/internal/service"
	"github.com/team4yf/fpm-iot-go-middleware/pkg/log"
	"github.com/team4yf/fpm-iot-go-middleware/pkg/pubsub"
	"github.com/team4yf/fpm-iot-go-middleware/router/middleware"
)

//App the core application
type App struct {
	Router     *mux.Router
	Middleware *middleware.Middleware
	PubSub     pubsub.PubSub
	Service    s.DeviceService
	m          alice.Chain
	LintaiAPI  rest.Client
}

func (app *App) Init() {

	app.Router = mux.NewRouter()
	app.PubSub = pubsub.NewMQTTPubSub(config.MqttConfig)
	app.Service = s.NewSimpleDeviceService()
	app.Middleware = &middleware.Middleware{}
	app.m = alice.New(app.Middleware.LoggerMiddleware, app.Middleware.RecoverMiddleware)

}

func (app *App) Run(addr string) {

	log.Infof("startup %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, app.Router))
}

func (app *App) Post(url string, handler func(w http.ResponseWriter, r *http.Request)) {
	app.Router.Handle(url, app.m.ThenFunc(handler)).Methods("POST")
}

func (app *App) Get(url string, handler func(w http.ResponseWriter, r *http.Request)) {
	app.Router.Handle(url, app.m.ThenFunc(handler)).Methods("GET")
}

func (app *App) WriteJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func (app *App) SendOk(w http.ResponseWriter, payload interface{}) {
	data, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (app *App) Fail(w http.ResponseWriter, result string) {
	err := errno.News(result)
	app.SendError(w, err)
}

func (app *App) FailWithError(w http.ResponseWriter, err error) {
	e := errno.NewsWithError(err)
	app.SendError(w, e)
}

func (app *App) FailWithCode(w http.ResponseWriter, code int, result string) {
	err := errno.NewsWithCode(code, result)
	app.SendError(w, err)
}

func (app *App) SendError(w http.ResponseWriter, err *errno.Errno) {
	data, _ := json.Marshal(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
