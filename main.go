package main

func init() {
}

func main() {
	cfg := &Config{}
	pubSub := cfg.GetPubSub()
	app := &App{}
	app.Init(pubSub)

	app.Run(":9000")
}
