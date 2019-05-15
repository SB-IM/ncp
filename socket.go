package main

import (
	"fmt"
  "net"
	"log"
  "strings"
  "time"
  "os"
)

func socketServer(addr string, input chan string, output chan string) {
  logger := log.New(os.Stdout, "[Socket Server] ", log.LstdFlags)
  ln, err := net.Listen("tcp", addr)
  if err != nil {
    logger.Println(err)
  }
  for {
    conn, err := ln.Accept()
    logger.Println("New connect")
    if err != nil {
      logger.Println(err)
    }
    go socketLink(conn, input, output)
  }
}

func socketClient(addr string, input chan string, output chan string) {
  logger := log.New(os.Stdout, "[Socket Client] ", log.LstdFlags)
  for {
    conn, err := net.Dial("tcp", addr)
    //defer conn.Close()
    if err != nil {
      //logger.Println(err)
      time.Sleep(1000000000)
      // handle error
    } else {
      logger.Println("New connect")
      socketLink(conn, input, output)
      logger.Println("Connect err try reconnect")
    }
  }
}

func socketLink(conn net.Conn, input chan string, output chan string) {
  go socketSend(conn, input)
  socketRecv(conn, output)
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

