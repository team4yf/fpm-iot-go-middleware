PROJECTNAME=$(shell basename "$(PWD)")
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin

define RECOMMENDED_TAG=
$(git tag -l | sort -V | tail -n 1 | sed 's#\.# #g' | awk '{ print $1"."$2 + 1".0" }')
endef

all: install build docker-build docker-push

install:
	go mod download

dev:
	go build -o $(GOBIN)/app ./main.go && $(GOBIN)/app

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(GOBIN)/app ./main.go

build-prod:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -tags prod -o $(GOBIN)/app ./main.go

beat:
	curl -H "Content-Type: application/json" -XPOST -d '{"data":"321123"}' localhost:9000/push/light/lt10/beat

sub:
	mosquitto_sub -h www.ruichen.top -t "^push/$(uuid)/event" -u "admin" -P "123123123"

create-redis-data:
	hset device:light:lt10:321123 projid 1
	hset device:light:lt10:321123 uuid d8f7r9fo

docker-build:
	docker build --tag fpm-iot-middleware:v2.0 --tag yfsoftcom/fpm-iot-middleware:v2.0 .

docker-push:
	docker push yfsoftcom/fpm-iot-middleware:v2.0

docker-run:
	docker run -e "REDIS_HOST=192.168.88.111" -p 9000:9000 fpm-iot-middleware:v2.0

docker-dev:
	docker-compose -f docker-compose.dev.yml up --build -d