
all: dep build

dep:
	sh get_version.sh

build:
	go build

run:
	go run `ls *.go | grep -v _test.go`

test:
	go test -cover

clean:
	go clean

