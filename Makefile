
GOPATH:=$(shell go env GOPATH)

.PHONY: proto test docker


build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w' -o test_v1 ./main.go

docker:
	docker build . -t test_v1:latest
