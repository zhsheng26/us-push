package main

import (
	"eusunpower.com/us-push/udp"
	"fmt"
	"github.com/nu7hatch/gouuid"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	port string = ":1200"
)

var pf = func(format string, a ...interface{}) {
	fmt.Printf(format, a)
}

type Server struct {
	conn     *net.UDPConn
	messages chan string
	clients  map[*uuid.UUID]Client
}

type Client struct {
	userID   uuid.UUID
	userName string
	userAddr *net.UDPAddr
}

type Message struct {
	messageType      udp.MessageType
	userID           *uuid.UUID
	userName         string
	content          string
	connectionStatus udp.ConnectionStatus
	time             string
}

func (server *Server) handleMessage() {
	var buf [512]byte

	n, addr, err := server.conn.ReadFromUDP(buf[0:])
	if err != nil {
		return
	}

	msg := string(buf[0:n])
	m := server.parseMessage(msg)

	if m.connectionStatus == udp.LEAVING {
		delete(server.clients, m.userID)
		server.messages <- msg
		pf("%s left", m.userName)
	} else {
		switch m.messageType {
		case udp.FUNC:
			var c Client
			c.userAddr = addr
			c.userID = *m.userID
			c.userName = m.userName
			server.clients[m.userID] = c
			server.messages <- msg
			pf("%s joining", m.userName)
		case udp.CLASSIQUE:
			pf("%s %s: %s", m.time, m.userName, m.content)
			server.messages <- msg
		}
	}
}

func (server *Server) parseMessage(msg string) (m Message) {
	stringArray := strings.Split(msg, "\x01")

	fmt.Println("")
	m.userID, _ = uuid.ParseHex(stringArray[0])
	messageTypeStr, _ := strconv.Atoi(stringArray[1])
	m.messageType = udp.MessageType(messageTypeStr)
	m.userName = stringArray[2]
	m.content = stringArray[3]
	m.time = stringArray[4]
	if strings.HasPrefix(msg, ":q") || strings.HasPrefix(msg, ":quit") {
		pf("%s is leaving", m.userName)
		m.connectionStatus = udp.LEAVING
	}
	return
}

func (server *Server) sendMessage() {
	for {
		msg := <-server.messages
		for _, c := range server.clients {
			_, err := server.conn.WriteToUDP([]byte(msg), c.userAddr)
			checkError(err)
		}
	}

}

func checkError(err error) {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Fatal error:%s", err.Error())
		os.Exit(1)
	}
}

func main() {
	udpAddress, err := net.ResolveUDPAddr("udp4", port)
	checkError(err)

	var s Server
	s.messages = make(chan string, 20)
	s.clients = make(map[*uuid.UUID]Client, 0)

	s.conn, err = net.ListenUDP("udp", udpAddress)
	checkError(err)

	go s.sendMessage()

	for {
		s.handleMessage()
	}
}
