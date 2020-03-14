GOPATH:=$(shell go env GOPATH)

.PHONY: proto
proto:
	protoc --proto_path=${GOPATH}/src:. --micro_out=paths=source_relative:. --go_out=paths=source_relative:. proto/*.proto
