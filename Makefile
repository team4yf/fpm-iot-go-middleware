
build:
	CGO_ENABLED=0 GOOS=linux go build -o app
beat: 
	curl -H "Content-Type: application/json" -XPOST -d '{"data":1}' localhost:9000/push/light/lb/beat
