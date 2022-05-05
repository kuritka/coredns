# Makefile for building CoreDNS
GITCOMMIT:=$(shell git describe --dirty --always)
BINARY:=coredns
SYSTEM:=
CHECKS:=check
BUILDOPTS:=-v
GOPATH?=$(HOME)/go
MAKEPWD:=$(dir $(realpath $(firstword $(MAKEFILE_LIST))))
CGO_ENABLED?=0

ORG_NAME=kuritka
PROVIDER_NAME=coredns



.PHONY: all
all: coredns

.PHONY: coredns
coredns: $(CHECKS)
	CGO_ENABLED=$(CGO_ENABLED) $(SYSTEM) go build $(BUILDOPTS) -ldflags="-s -w -X github.com/coredns/coredns/coremain.GitCommit=$(GITCOMMIT)" -o $(BINARY)

.PHONY: check
check: core/plugin/zplugin.go core/dnsserver/zdirectives.go

core/plugin/zplugin.go core/dnsserver/zdirectives.go: plugin.cfg
	go generate coredns.go
	go get

.PHONY: gen
gen:
	go generate coredns.go
	go get

.PHONY: pb
pb:
	$(MAKE) -C pb

.PHONY: clean
clean:
	go clean
	rm -f coredns

.PHONY: build
build:
	go build -o ./coredns

.PHONY: run
run: build
	./coredns -conf Corefile -p 5053

image: build
	docker build . -t $(ORG_NAME)/$(PROVIDER_NAME):$(VERSION) -f Dockerfile

image-push:
	# docker push $(ORG_NAME)/$(PROVIDER_NAME):$(VERSION)
	k3d image import $(ORG_NAME)/$(PROVIDER_NAME):$(VERSION) -c abc

apply:
	#kubectl delete -f infrastructure/infra.yaml
	kubectl apply -f infrastructure/infra.yaml

deploy: build image image-push apply

push-remote: build image
	docker push $(ORG_NAME)/$(PROVIDER_NAME):$(VERSION)

VERSION=v0.0.4
