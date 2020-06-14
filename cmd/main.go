package main

func init() {
}

func main() {
	cfg := &Config{}
	pubSub := cfg.GetPubSub()
	service := cfg.GetService()
	app := &App{}
	app.Config = cfg
	app.Init(pubSub, service)

	app.Run(cfg.GetConfigOrDefault("server.host", "0.0.0.0") + ":" + cfg.GetConfigOrDefault("server.port", "9000"))
}
