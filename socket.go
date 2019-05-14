package main

import (
	"fmt"
  "net"
	"log"
  "strings"
  "time"
  "os"
)

func socketListen(addr string, input chan string, output chan string) {
  ln, err := net.Listen("tcp", addr)
  if err != nil {
    // handle error
    log.Println(err)
  }
  for {
    conn, err := ln.Accept()
    if err != nil {
      log.Println(err)
      // handle error
    }
    go socketLink(conn, input, output)
  }

}

func socketLink(conn net.Conn, input chan string, output chan string) {
  go socketSend(conn, input)
  socketRecv(conn, output)
}

func socketClient(host string, input chan string, output chan string) {
  for {
    conn, err := net.Dial("tcp", host)
    //defer conn.Close()
    if err != nil {
      log.Println(err)
      time.Sleep(1000000000)
      // handle error
    } else {
      socketLink(conn, input, output)
    }
  }
}

func socketRecv(conn net.Conn, ch chan string) {
  logger := log.New(os.Stdout, "[Socket recv] ", log.LstdFlags)
  buf := make([]byte, 4096)
  for {
    cnt, err := conn.Read(buf)
    if err != nil || cnt == 0 {
      fmt.Println("EEEEEEEEEEee")
      conn.Close()
      break
    }
    inStr := strings.TrimSpace(string(buf[0:cnt]))
    logger.Println(inStr)
    ch <- inStr
  }
}

func socketSend(conn net.Conn, ch chan string) {
  logger := log.New(os.Stdout, "[Socket send] ", log.LstdFlags)
  for {
    m := <- ch
    logger.Println(m)
    _, e := conn.Write([]byte(m + "\n"))

    if e != nil {
      log.Printf("Error: ")
      log.Println(e)
      ch <- m
      break
    }
  }
}


