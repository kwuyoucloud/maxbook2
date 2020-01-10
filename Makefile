# Go parameters
VERSION=1.12.9
GOCMD=go
GORUN=$(GOCMD) run
GOBUILD=$(GOCMD) build
GOFILE=./maxbook
MAINGOFILE=cmd/main.go
MAINGO=main

all: build

run:
	@$(GORUN) $(MAINGOFILE)
build:
	@$(GOBUILD) $(MAINGOFILE)
	@mv $(MAINGO) $(GOFILE)
	@echo "Build correct."

clean:
	@rm $(GOFILE)
