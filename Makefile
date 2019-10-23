
all: dep build

dep:
	cp -r golib/starlight/light .
	cp -r golib/starlight/gstreamer-src .
	sed -i 's#gst "gstreamer-src"#gst "./gstreamer-src"#g' light/gst2webrtc.go
	sed -i 's#"light"#"./light"#g' http.go
	sed -i 's#"light"#"./light"#g' ncpcmd.go
	go mod vendor
	sed -i 's#"./light"#"light"#g' ncpcmd.go
	sed -i 's#"./light"#"light"#g' http.go
	sed -i 's#gst "./gstreamer-src"#gst "gstreamer-src"#g' light/gst2webrtc.go
	mv gstreamer-src vendor
	mv light vendor

build:
	go build -mod=vendor

run:
	go run -mod=vendor `ls *.go | grep -v _test.go`

test:
	go test -mod=vendor -cover

