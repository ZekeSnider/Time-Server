GOPATH := $(CURDIR)
all: build

build:

	go fmt src/command/TimeServer/TimeServer.go
	GOPATH=$(GOPATH) go build -o bin/AuthServer src/command/config/config.go
	GOPATH=$(GOPATH) go build -o bin/AuthServer src/command/authserver/AuthServer.go
	GOPATH=$(GOPATH) go build -o bin/TimeServer src/command/timeserver/TimeServer.go

run:
	go fmt src/TimeServer.go
	GOPATH=$(GOPATH) go run src/TimeServer.go
