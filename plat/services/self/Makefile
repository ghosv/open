srv_name:="srv-self"

docker_name:=open-${srv_name}
GOPATH:=$(shell go env GOPATH)

.PHONY: build
build: proto
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${srv_name} main.go plugin.go

.PHONY: proto
proto:
	protoc --proto_path=${GOPATH}/src:../../..:. --micro_out=paths=source_relative:. --go_out=paths=source_relative:. proto/user.proto
	protoc --proto_path=${GOPATH}/src:../../..:. --micro_out=paths=source_relative:. --go_out=paths=source_relative:. proto/app.proto
	protoc --proto_path=${GOPATH}/src:../../..:. --micro_out=paths=source_relative:. --go_out=paths=source_relative:. proto/org.proto
	protoc --proto_path=${GOPATH}/src:../../..:. --micro_out=paths=source_relative:. --go_out=paths=source_relative:. proto/group.proto
	protoc --proto_path=${GOPATH}/src:../../..:. --micro_out=paths=source_relative:. --go_out=paths=source_relative:. proto/0-service.proto

.PHONY: dev
dev: proto
	go run . --registry consul --registry_address 127.0.0.1:8500

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t ${docker_name}:latest
