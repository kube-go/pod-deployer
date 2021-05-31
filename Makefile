OS ?= $(shell uname -s | tr A-Z a-z)
ARCH ?= $(shell go env GOARCH)

build-local:
	GOOS=$(OS) GOARCH=$(ARCH) CGO_ENABLED=0 go build -o pod-deployer
