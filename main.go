package main

func init() {
}

func main() {
	cfg := &Config{}
	pubSub := cfg.GetPubSub()
	service := cfg.GetService()
	app := &App{}
	app.Init(pubSub, service)

	app.Run(":9000")
}
