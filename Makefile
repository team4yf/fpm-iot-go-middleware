PROJECTNAME=$(shell basename "$(PWD)")
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin

define RECOMMENDED_TAG=
$(git tag -l | sort -V | tail -n 1 | sed 's#\.# #g' | awk '{ print $1"."$2 + 1".0" }')
endef

all: install build docker-build docker-push

local: install build docker-build-local docker-push
install:
	go mod download

dev:
	go build -o $(GOBIN)/app ./main.go && $(GOBIN)/app

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(GOBIN)/app ./main.go

build-prod:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -tags prod -o $(GOBIN)/app ./main.go
sub:
	mosquitto_sub -h open.yunplus.io -t '$d2s/ceaa191a/partner/push' -u "fpmuser" -P

pub:
	mosquitto_pub -h open.yunplus.io -t '$d2s/ceaa191a/partner/push' -u "fpmuser" -m "test"

docker-build-local:
	docker build -f Dockerfile-local --tag fpm-iot-middleware:v2.0 --tag yfsoftcom/fpm-iot-middleware:v2.0 .

docker-build:
	docker build --tag fpm-iot-middleware:v2.0 --tag yfsoftcom/fpm-iot-middleware:v2.0 .

docker-push:
	docker push yfsoftcom/fpm-iot-middleware:v2.0