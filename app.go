package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

type App struct {
	Router     *mux.Router
	Middleware *Middleware
	pubSub     PubSub
	service    Service
}

func (app *App) Init(pubSub PubSub, service Service) {
	app.Router = mux.NewRouter()
	app.pubSub = pubSub
	app.service = service
	app.Middleware = &Middleware{}
	m := alice.New(app.Middleware.LoggerMiddleware, app.Middleware.RecoverMiddleware)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	app.Router.Handle("/push/{device}/{brand}/{event}", m.ThenFunc(app.pushHandler)).Methods("POST")
}

func (app *App) Run(addr string) {
	log.Printf("startup %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, app.Router))
}

func (app *App) pushHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	device := params["device"]
	brand := params["brand"]
	event := params["event"]
	body, err := GetBodyString(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	go func() {
		var deviceID string
		if res, err := GetJsonPathData(body, "$.data"); err != nil {
			log.Printf("device id not: %v", err)
			return
		} else {
			deviceID = res.(string)
		}
		if uuid, err := app.service.Receive(device, brand, event, deviceID); err != nil {
			log.Printf("Error reading body: %v", err)
			return
		} else {
			app.pubSub.Publish(fmt.Sprintf("^push/%s/event", uuid), body)
		}
	}()

	log.Printf("device: %s brand:%s event:%s body:%s\n", device, brand, event, body)
	writeJSON(w, 200, (`{"errorNo":"ok"`))
}

func writeJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
