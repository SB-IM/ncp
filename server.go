package main

import (
	"log"
	"net"
	"strings"
	"sync"
)

type Link struct {
	conn    *net.Conn
	methods []string
}

type SocketServer struct {
	links  []*Link
	lock   sync.Mutex
	logger *log.Logger
}

func (this *SocketServer) Listen(addr string, input chan string, output chan string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		this.logger.Println("Listener:", addr, err)
	}
	defer listener.Close()

	go this.send(input)

	for {
		conn, err := listener.Accept()
		this.logger.Println("New Connect:", &conn)
		if err != nil {
			this.logger.Println("Connect Err:", &conn, err)
		} else {
			go this.recv(&conn, output)
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
			this.logger.Println("Send:", conn, msg)
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
			this.logger.Println("Connect Err:", conn, err)
			break
		}
		msg := strings.TrimSpace(string(buf[0:cnt]))
		methods, result := getReg([]byte(msg))
		if len(methods) != 0 {
			link.methods = methods
			this.addLink(link)
			this.logger.Println("Method Reg:", conn, methods)
			if result != "" {
				(*conn).Write([]byte(result + "\n"))
			}
		} else {
			this.logger.Println("Recv:", conn, msg)
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
