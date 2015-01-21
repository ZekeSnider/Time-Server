all: build

build:
	go build -o bin/TimeServer TimeServer.go

run:
	go run TimeServer.go
