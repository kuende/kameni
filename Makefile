.PHONY: default release deps

default:
	mkdir -p bin
	go build -o bin/kameni

release:
	mkdir -p release
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w' -o release/kameni-v0.1.1-linux-amd64 .
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -tags netgo -ldflags '-w' -o release/kameni-v0.1.1-darwin-amd64 .

deps:
	godep save -r ./...
