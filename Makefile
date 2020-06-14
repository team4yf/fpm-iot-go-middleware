
dev:
	go build -o app && ./app
build:
	CGO_ENABLED=0 GOOS=linux go build -o app
beat: 
	curl -H "Content-Type: application/json" -XPOST -d '{"data":"321123"}' localhost:9000/push/light/lt10/beat
sub:
	mosquitto_sub -h www.ruichen.top -t "^push/$(uuid)/event" -u "admin" -P "123123123"
create-redis-data:
	hset device:light:lt10 321123 d8f7r9fo

docker-build:
	docker build --tag fpm-iot-middleware:v2.0 .

docker-run:
	docker run -e "REDIS_HOST=192.168.88.111" -p 9000:9000 fpm-iot-middleware:v2.0

docker-dev:
	docker-compose -f docker-compose.dev.yml up --build -d