GOPATH := $(CURDIR)
all: build

build:
	GOPATH=$(GOPATH) go fmt src/command/timeserver/timeserver.go
	GOPATH=$(GOPATH) go fmt src/command/authserver/authserver.go
	GOPATH=$(GOPATH) go install command/config
	GOPATH=$(GOPATH) go build -o bin/AuthServer src/command/authserver/authserver.go
	GOPATH=$(GOPATH) go build -o bin/TimeServer src/command/timeserver/timeserver.go

run:
	go fmt src/TimeServer.go
	GOPATH=$(GOPATH) go run src/command/authserver/authserver.go
	GOPATH=$(GOPATH) go run src/command/timeserver/timeserver.go
 
install:
	GOPATH=$(GOPATH) go install command/timeserver
	GOPATH=$(GOPATH) go install command/authserver
