package main

import (
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

type Link struct {
	conn    *net.Conn
	methods []string
}

type SocketServer struct {
	links []*Link
	lock  sync.Mutex
}

func (this *SocketServer) Listen(addr string, file *os.File, input chan string, output chan string) {
	logger := log.New(file, "[Server] ", log.LstdFlags)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Println(err)
	}
	defer listener.Close()

	go this.send(input)

	for {
		conn, err := listener.Accept()
		logger.Println("New connect")
		if err != nil {
			logger.Println(err)
		} else {
			this.links = append(this.links, &Link{conn: &conn})
			logger.Println(this.links)
			go this.recv(&conn, input)
		}
	}
}

func (this *SocketServer) addLink(link *Link) {
	this.delLink(link)
	this.lock.Lock()
	this.links = append(this.links, link)
	this.lock.Unlock()
}

func (this *SocketServer) delLink(link *Link) {
	this.lock.Lock()
	for index, run_link := range this.links {
		if run_link == link {

			// Order is not important
			this.links[index] = this.links[len(this.links)-1]
			this.links = this.links[:len(this.links)-1]
		}
	}
	this.lock.Unlock()
}

func (this *SocketServer) send(input chan string) {
	for msg := range input {
		for _, conn := range this.getMethodMatchConns(getJSONRPC(msg).Method) {
			(*conn).Write([]byte(msg + "\n"))
		}
	}
}

func (this *SocketServer) recv(conn *net.Conn, output chan string) {
	link := &Link{
		conn: conn,
	}

	defer func() {
		this.delLink(link)
		(*conn).Close()
	}()

	buf := make([]byte, 4096)
	for {
		cnt, err := (*conn).Read(buf)
		if err != nil || cnt == 0 {
			break
		}
		msg := strings.TrimSpace(string(buf[0:cnt]))
		methods, ok := getReg([]byte(msg))
		if ok {
			link.methods = methods
			this.addLink(link)
		} else {
			output <- msg
		}
	}
}

func (this *SocketServer) getMethodMatchConns(commond string) []*net.Conn {
	var conns []*net.Conn
	for _, link := range this.links {
		for _, method := range link.methods {
			if commond == method {
				conns = append(conns, link.conn)
			}
		}
	}
	return conns
}

func (this *SocketServer) GetMethods() []string {
	var methods []string
	for _, link := range this.links {
		methods = append(methods, link.methods...)
	}
	return methods
}
