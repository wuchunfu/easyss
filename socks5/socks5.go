package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

const (
	SERVER_PORT = 1080
)

const (
	READ_DEADLINE  = time.Minute
	WRITE_DEADLINE = 2 * time.Minute
)

func handleRequest(conn net.Conn) {
	if conn == nil {
		return
	}
	defer conn.Close()
	conn.SetReadDeadline(time.Now().Add(READ_DEADLINE))
	conn.SetWriteDeadline(time.Now().Add(WRITE_DEADLINE))

	var b [1024]byte
	if _, err := conn.Read(b[:]); err != nil {
		log.Println("conn read err:", err.Error())
		return
	}

	// only handle socks5 protocol
	if b[0] != 0x05 {
		log.Println("server do not support client version:", b[0])
		return
	}

	if _, err := conn.Write([]byte{0x05, 0x00}); err != nil {
		log.Println("client cannot arrived, err:", err.Error())
		return
	}

	n, err := conn.Read(b[:])
	if err != nil {
		log.Println("conn read err:", err.Error())
		return
	}
	if b[0] != 0x05 || b[1] != 0x01 || b[2] != 0x00 {
		log.Println("client cmd param is not supported:", b[0], b[1], b[2])
		return
	}

	var targetHost, targetPort string
	switch b[3] {
	case 0x01: // ipv4
		targetHost = net.IPv4(b[4], b[5], b[6], b[7]).String()
	case 0x03:
		targetHost = string(b[5 : n-2])
	case 0x04:
		targetHost = net.IP(b[4:20]).String()
	}
	targetPort = strconv.Itoa(int(b[n-2])<<8 | int(b[n-1]))

	targetConn, err := net.Dial("tcp", net.JoinHostPort(targetHost, targetPort))
	if err != nil {
		log.Printf("net dial host:%v, port:%v, err:%v\n", targetHost, targetPort, err)
		return
	}
	defer targetConn.Close()

	if _, err := conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}); err != nil {
		log.Println("client cannot arrived, err:", err.Error())
		return
	}

	go io.Copy(targetConn, conn)
	io.Copy(conn, targetConn)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", SERVER_PORT))
	if err != nil {
		log.Panic(err)
	}
	log.Println("socks5 server listen on:", SERVER_PORT)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("accept conn err:", err.Error())
			continue
		}

		go handleRequest(conn)
	}
}
