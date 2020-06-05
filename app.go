package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	pubSub PubSub
}

func (app *App) Init(pubSub PubSub) {
	app.Router = mux.NewRouter()
	app.pubSub = pubSub
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	app.Router.HandleFunc("/push/{device}/{brand}/{event}", app.pushHandler).Methods("POST")
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
	var body interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		panic(nil)
	}
	defer r.Body.Close()
	go func() {
		app.pubSub.Publish("$push/event/dd", body)
	}()

	log.Printf("device: %s brand:%s event:%s body:%s\n", device, brand, event, body)
	writeJSON(w, 200, `{"errorNo":"ok"`)
}

func writeJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
