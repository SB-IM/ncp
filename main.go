package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"time"

	C "sb.im/ncp/constant"
	"sb.im/ncp/ncpio"

	cfg "sb.im/ncp/tests/help"
)

func main() {
	config_path := "config.yml"

	help := flag.Bool("h", false, "this help")
	debug := flag.Bool("debug", false, "use debug mode")
	flag.StringVar(&config_path, "c", "config.yml", "set configuration file")

	show_version := flag.Bool("v", false, "show version")
	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	if *show_version {
		fmt.Printf("Ncp %s %s %s %s %s\n", C.Version, C.Commit, runtime.GOOS, runtime.GOARCH, C.BuildTime)
		return
	}

	if os.Getenv("NCP_CONF") != "" {
		config_path = os.Getenv("NCP_CONF")
	}
	log.Println("load config: " + config_path)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	config, err := cfg.GetConfig(config_path)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("%+v", config)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ncp := ncpio.NewNcpIOs(config.NcpIO)
	if *debug {
		ncp.Debuger = log.New(os.Stdout, "[NCPIO] ", log.LstdFlags)
	}
	go ncp.Run(ctx)

	// Wait mqttd server startup && sub topic on broker
	time.Sleep(3 * time.Millisecond)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	<-sigs
	log.Println("ncpio exit")
}
