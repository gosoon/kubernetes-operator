GOFILES := $(shell find . -name '*.go' | grep -v -E '(./vendor)')
BIN := $(shell basename $(CURDIR))
# PLATFORM := $(shell go env GOOS)
PLATFORM := linux
AnsibleinitPath := cmd/ansible/ansibleinit.go
OperatorPath := cmd/kube-operator/app.go

# ifneq (PLATFORM, "windows")
# 	PLATFORM = linux
# endif

.DEFAULT_GOAL := default

default: $(PLATFORM)

all: linux darwin windows

linux: output/$(BIN)
darwin: output/$(BIN)
windows: output/$(BIN)

images: 
	mkdir -p output
	GO111MODULE=on GOPROXY=https://mirrors.aliyun.com/goproxy/ GOOS=$(PLATFORM) GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o output/ansibleinit $(AnsibleinitPath)
	docker build --no-cache -f build/ansible/Dockerfile -t ansibleinit .
	
	GO111MODULE=on GOPROXY=https://mirrors.aliyun.com/goproxy/ GOOS=$(PLATFORM) GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o output/kubernetes-operator $(OperatorPath)
	docker build --no-cache -f build/kube-operator/Dockerfile -t kubernetes-operator .
	@rm -rf output 

check:
	@find . -name vendor -prune -o -name '*.go' -exec gofmt -s -d {} +
	@go vet $(shell go list ./... | grep -v '/vendor/')
	@go test -v $(shell go list ./... | grep -v '/vendor/')

vendor:
	dep ensure

clean:
	@rm -rf output

output/%: LDFLAGS=-s -w
output/%: $(GOFILES)
	mkdir -p $(dir $@)
	GO111MODULE=on GOPROXY=https://mirrors.aliyun.com/goproxy/ GOOS=$(PLATFORM) GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $@
