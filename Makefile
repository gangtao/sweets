
ROOT := $(CURDIR)
BIN_NAME ?= sweets
VERSION ?= 0.1
IMAGE_NAME ?= $(BIN_NAME):$(VERSION)
DOCKER_ID_USER ?= naughtytao

FULLNAME=$(DOCKER_ID_USER)/${BIN_NAME}:${VERSION}

build:
	cd ./server ; env GOOS=linux GOARCH=amd64 go build
	cd ./client ; env GOOS=linux GOARCH=amd64 go build

docker: Dockerfile build
	docker build -t $(IMAGE_NAME) .

clean:
	rm -f ./server/server
	rm -f ./client/client
	rm -f ./zk/zk
	rm -f ./etcd/etcd

push:
	docker tag $(IMAGE_NAME) ${FULLNAME}
	docker push ${FULLNAME}
