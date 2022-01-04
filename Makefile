ifeq ($(GOPATH),)
GOPATH := $(HOME)/go
endif

all: lint vet build

build: server client

server:
	@cd cmd/$@ && go build -o ../../bin/$@ -tags osusergo,netgo
	
client:
	@cd cmd/$@ && go build -o ../../bin/$@ -tags osusergo,netgo

vet:
	@go vet ./...

lint:
	@revive ./...

clean:
	@rm -rf bin
