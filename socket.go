package main

import (
	"fmt"
  "net"
	"log"
  "strings"
  "time"
  "os"
  "regexp"

  mqtt "github.com/eclipse/paho.mqtt.golang"
)

func socketServerTran(addr string, client mqtt.Client, topic_prefix string) {
  logger := log.New(os.Stdout, "[Socket Tran] ", log.LstdFlags)
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

    go socketTran(conn, client, topic_prefix)
  }
}

func socketTran(conn net.Conn, client mqtt.Client, topic_prefix string) {
  logger := log.New(os.Stdout, "[Socket Tran] ", log.LstdFlags)
  ch := make(chan string, 100)
  go func() {
    tag := "default"
    for x := range ch {
      matched, _ := regexp.MatchString(`^\-tran\:`, x)
      if matched {
        tag = strings.Split(strings.Split(x, " ")[1], "/")[2]
        logger.Println("Set tag: " + tag)
      } else {
        client.Publish(topic_prefix + "/msg/" + tag, 0, true, x)
      }
    }
  }()

  go func() {
    buf := make([]byte, 4096)
    for {
      cnt, err := conn.Read(buf)
      if err != nil || cnt == 0 {
        logger.Println("Error: socket close")
        conn.Close()
        close(ch)
        break
      }
      inStr := strings.TrimSpace(string(buf[0:cnt]))
      logger.Println(inStr)
      ch <- inStr
    }
  }()
}

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

