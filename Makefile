.PHONY: all
all: build

.PHONY: build
build:
	go build -gcflags=all=-d=checkptr=1 -trimpath -o bin/meilidex-upload main.go

build-linux:
	GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -trimpath -o bin/meilidex-upload main.go
