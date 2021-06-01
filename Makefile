#!/bin/bash

export APPNSP=clcng
export APPNAME=bitcoin-wallet
export APPVERSION=0.0.1
export REGISTRY=github.com
export DISTDIR=dist
export GOOS=linux
export GOARCH=amd64

SRCBRANCH := $(shell git branch | grep \* | cut -d ' ' -f2)
ifeq ($(SRCBRANCH),$(filter lab uat prd, $(SRCBRANCH)))
	SRCRELEASE:=release-$(SRCBRANCH)
else
	SRCRELEASE=latest	
endif

APPRELEASE ?= $(SRCRELEASE)
APPBUILD ?= $(shell date +'%Y%m%d%H%M%S')-$(shell git rev-parse --short=12 HEAD)

.PHONY: build-fast build

build-fast:
	mkdir -p $(DISTDIR)
	env CGO_ENABLED=0 env GOOS=$(GOOS) env GOARCH=$(GOARCH) go build -ldflags="-s -w -X main.version=$(APPVERSION) -X main.build=$(APPBUILD)" -o "$(DISTDIR)/$(GOOS)-$(GOARCH)/$(APPNAME)"
	docker build \
	-f Dockerfile.fast \
	-t $(REGISTRY)/$(APPNSP)/$(APPNAME):$(APPVERSION) .

build:
	docker build \
	--build-arg APPVERSION=$(APPVERSION) \
	--build-arg APPBUILD=$(APPBUILD) \
	-t $(REGISTRY)/$(APPNSP)/$(APPNAME):$(APPVERSION) .
