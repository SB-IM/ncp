OS=
ARCH=
NAME=ncp
BINDIR=bin
PROFIX=
VERSION=$(shell git describe --tags || echo "unknown version")
BUILDTIME=$(shell date)
GOBUILD=go build -ldflags '-X "sb.im/ncp/constant.Version=$(VERSION)" \
				-X "sb.im/ncp/constant.BuildTime=$(BUILDTIME)"'

PLATFORM_LIST = \
								darwin-amd64 \
								linux-386 \
								linux-amd64 \
								linux-armv7 \
								linux-armv8 \
								freebsd-amd64

WINDOWS_ARCH_LIST = \
										windows-386 \
										windows-amd64

all: build

build:
	$(GOBUILD)

run:
	go run `ls *.go | grep -v _test.go` -debug

install:
	install -Dm755 ncp -t ${PROFIX}/usr/bin/
	install -Dm644 conf/ncp.service -t ${PROFIX}/lib/systemd/system/
	install -Dm644 conf/ncp@.service -t ${PROFIX}/lib/systemd/system/
	install -Dm644 conf/config-dist.yml -t ${PROFIX}/etc/ncp/

# Need Container Network Interface
# Linux tc (Traffic Control)
#
# Manual test
# docker run --cap-add "NET_ADMIN" -it -v $(pwd):/ncp  golang:1.13.1-buster /bin/bash
# apt-get update -y && apt-get install -y mosquitto-clients
# cd /ncp
#
# docker run eclipse-mosquitto:1.6
#
# # YOU Broker IP
# MQTT=172.17.0.3:1883 ./test.network
test-detach:
	CGO_ENABLED=0 go test ./tests/network -c -o test.network -v

# Need mosquitto && mosquitto_pub
test-integration:
	go test ./tests/integration

test:
	go test ./ncpio ./util ./history ./cache -cover

# \(statements\)(?:\s+)?(\d+(?:\.\d+)?%)
# https://stackoverflow.com/questions/61246686/go-coverage-over-multiple-package-and-gitlab-coverage-badge
cover:
	go test ./ncpio ./util ./history -coverprofile profile.cov
	go tool cover -func profile.cov
	@rm profile.cov

darwin-amd64:
	GOARCH=amd64 GOOS=darwin $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

linux-386:
	GOARCH=386 GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

linux-amd64:
	GOARCH=amd64 GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

linux-armv7:
	GOARCH=arm GOOS=linux GOARM=7 $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

linux-armv8:
	GOARCH=arm64 GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

freebsd-amd64:
	GOARCH=amd64 GOOS=freebsd $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

windows-386:
	GOARCH=386 GOOS=windows $(GOBUILD) -o $(BINDIR)/$(NAME)-$@.exe

windows-amd64:
	GOARCH=amd64 GOOS=windows $(GOBUILD) -o $(BINDIR)/$(NAME)-$@.exe

releases: $(PLATFORM_LIST) $(WINDOWS_ARCH_LIST)

clean:
	go clean
	rm $(BINDIR)/*

