
dev:
	go build -o app && ./app
build:
	CGO_ENABLED=0 GOOS=linux go build -o app
beat: 
	curl -H "Content-Type: application/json" -XPOST -d '{"data":"321123"}' localhost:9000/push/light/lt/beat
sub:
	mosquitto_sub -h www.ruichen.top -t "^push/$(uuid)/event" -u "admin" -P "123123123"
create-redis-data:
	hset device:light:lt 321123 d8f7r9fo