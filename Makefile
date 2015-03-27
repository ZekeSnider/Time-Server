GOPATH := $(CURDIR)
all: build

build:
	GOPATH=$(GOPATH) go fmt src/command/timeserver/timeserver.go
	GOPATH=$(GOPATH) go fmt src/command/authserver/authserver.go
	GOPATH=$(GOPATH) go fmt src/command/loadgen/loadgen.go
	GOPATH=$(GOPATH) go fmt src/command/monitor/monitor.go
	GOPATH=$(GOPATH) go build command/config
	GOPATH=$(GOPATH) go build command/counter
	GOPATH=$(GOPATH) go build -o bin/authserver src/command/authserver/authserver.go
	GOPATH=$(GOPATH) go build -o bin/timeserver src/command/timeserver/timeserver.go
	GOPATH=$(GOPATH) go build -o bin/loadgen src/command/loadgen/loadgen.go
	GOPATH=$(GOPATH) go build -o bin/monitor src/command/monitor/monitor.go

run:
	go fmt src/TimeServer.go
	GOPATH=$(GOPATH) go run src/command/authserver/authserver.go
	GOPATH=$(GOPATH) go run src/command/timeserver/timeserver.go
 
install:
	GOPATH=$(GOPATH) go install command/timeserver
	GOPATH=$(GOPATH) go install command/authserver

time:
	GOPATH=$(GOPATH) go build -o bin/timeserver src/command/timeserver/timeserver.go

auth:
	GOPATH=$(GOPATH) go build -o bin/authserver src/command/authserver/authserver.go

monitor:
	GOPATH=$(GOPATH) go build -o bin/monitor src/command/monitor/monitor.go

load:
	GOPATH=$(GOPATH) go build -o bin/LoadGen src/command/loadgen/loadgen.go

counter:
	GOPATH=$(GOPATH) go build src/command/counter/counter.go

loadtest:
	./bin/authserver -log=auth-log.xml &
	./bin/timeserver -log=seelog.xml -port=8081 -maxinflight=80 -response=500 -deviation=300 &
	./bin/loadgen -url='http://localhost:8081/time' -runtime=10 -rate=200 --burst=20 -timeout=1000