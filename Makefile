GOPATH := $(CURDIR)
all: build

build:

	go fmt src/TimeServer.go
	GOPATH=$(GOPATH) go build -o bin/TimeServer src/TimeServer.go

run:
	go fmt src/TimeServer.go
	GOPATH=$(GOPATH) go run src/TimeServer.go
