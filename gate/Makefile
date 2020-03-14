srv_name:="gateway"

docker_name:=open-${srv_name}
GOPATH:=$(shell go env GOPATH)

.PHONY: build
build: schema
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${srv_name} main.go plugin.go

.PHONY: dev
dev: schema
	go run . --registry consul --registry_address 127.0.0.1:8500

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t ${docker_name}:latest


.PHONY: schema
schema:
	cd graphql/schema; go generate
