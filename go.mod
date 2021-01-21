module sb.im/ncp

go 1.13

require (
	github.com/SB-IM/jsonrpc-lite v0.1.0
	github.com/eclipse/paho.golang v0.9.1-0.20210104204216-f20508bef4fd
	github.com/fastly/go-utils v0.0.0-20180712184237-d95a45783239 // indirect
	github.com/jehiah/go-strftime v0.0.0-20171201141054-1d33003b3869 // indirect
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/lestrrat-go/strftime v1.0.3 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/tebeka/strftime v0.1.5 // indirect
	gopkg.in/yaml.v2 v2.2.7
)

replace github.com/eclipse/paho.golang v0.9.1-0.20210104204216-f20508bef4fd => github.com/SB-IM/paho.golang v0.9.1-0.20210121044738-926a17219d78
