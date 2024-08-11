PROJECTNAME=$(shell basename "$(PWD)")
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/cmd
GOFILES=$(wildcard *.go)
BUILDDATE=$(shell date +'%Y/%m/%d')
BUILDVERSION=$(shell git rev-parse --short HEAD)
CRYPTOKEYS=$(GOBASE)/internal/crypto
SAVEFILES=$(GOBASE)/static/client

install: go-get go-build

go-get:
	@echo "  >  Checking if there is any missing dependencies..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go get ./...

go-build:
	@echo "  >  Building binary..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go build -ldflags "-X 'main.buildDate=$(BUILDDATE)' -X 'main.buildVersion=$(BUILDVERSION)' -X 'main.buildCryptoKeysPath=$(CRYPTOKEYS)' -X 'main.buildSaveFilePath=$(SAVEFILES)'" -o $(PROJECTNAME) $(GOBIN)/client $(GOFILES)

go-run:
	@echo " > run program..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN)
