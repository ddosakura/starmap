srv_name:="auth-srv"

docker_name:=starmap-${srv_name}
#GOPATH:=$(shell go env GOPATH)
GOPATH:=${shell echo $$ORIGIN_GOPATH}

.PHONY: build
build: proto
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${srv_name} main.go plugin.go

.PHONY: proto
proto:
	protoc --proto_path=${GOPATH}/src:. --micro_out=. --go_out=. proto/*.proto

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t ${docker_name}:latest
