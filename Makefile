VERSION=$(shell git describe --tags || echo "unknown version")
BUILDTIME=$(shell date -u)
GOBUILD=go build -ldflags '-X "sb.im/ncp/constant.Version=$(VERSION)" \
				-X "sb.im/ncp/constant.BuildTime=$(BUILDTIME)"'

all: build

build:
	$(GOBUILD)

run:
	go run `ls *.go | grep -v _test.go`

test:
	go test -cover

clean:
	go clean

