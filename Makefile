.PHONY: default release deps

default:
	mkdir -p bin
	go build -o bin/kameni

release:
	mkdir -p release
	CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' -o release/kameni .

deps:
	godep save -r ./...
