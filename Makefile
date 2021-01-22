OS=
ARCH=
PROFIX=
VERSION=$(shell git describe --tags || echo "unknown version")
BUILDTIME=$(shell date)
GOBUILD=GOOS=$(OS) GOARCH=$(ARCH) \
				go build -ldflags '-X "sb.im/ncp/constant.Version=$(VERSION)" \
				-X "sb.im/ncp/constant.BuildTime=$(BUILDTIME)"'

all: build

build:
	$(GOBUILD)

run:
	go run `ls *.go | grep -v _test.go`

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
	go test ./ncpio ./util ./history -cover

# \(statements\)(?:\s+)?(\d+(?:\.\d+)?%)
# https://stackoverflow.com/questions/61246686/go-coverage-over-multiple-package-and-gitlab-coverage-badge
cover:
	go test ./ncpio ./util ./history -coverprofile profile.cov
	go tool cover -func profile.cov
	@rm profile.cov

clean:
	go clean

