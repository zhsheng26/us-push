package main

import (
	"bufio"
	"eusunpower.com/us-push/udp"
	"fmt"
	"github.com/nu7hatch/gouuid"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	connection          *net.UDPConn
	alive               bool
	userID              uuid.UUID
	userName            string
	sendingMessageQueue chan string
	receiveMessages     chan string
}

func (c *Client) packMessage(msg string, messageType udp.MessageType) string {
	return strings.Join([]string{c.userID.String(), strconv.Itoa(int(messageType)), c.userName, msg, time.Now().Format("15:04:05")}, "\x01")
}

func (c *Client) funcSendMessage(msg string) {
	message := c.packMessage(msg, udp.FUNC)
	_, err := c.connection.Write([]byte(message))
	checkError(err, "func_sendMessage")
}

func (c *Client) sendMessage() {
	for c.alive {
		msg := <-c.sendingMessageQueue
		message := c.packMessage(msg, udp.CLASSIQUE)
		_, err := c.connection.Write([]byte(message))
		checkError(err, "sendMessage")
	}

}

func (c *Client) receiveMessage() {
	var buf [512]byte
	//var userID *uuid.UUID
	for c.alive {
		n, err := c.connection.Read(buf[0:])
		checkError(err, "receiveMessage")
		c.receiveMessages <- string(buf[0:n])
		fmt.Println("")
	}
}

func (c *Client) readInput() {
	var msg string
	for c.alive {
		fmt.Println("msg: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			msg = scanner.Text()
			if msg == ":quit" || msg == ":q" {
				c.alive = false
			} else if msg == ":clear" {
				udp.CallClear()
				msg = "udp"
			}
			c.sendingMessageQueue <- msg
		}
	}
}

func (c *Client) printMessage() {
	for c.alive {
		msg := <-c.receiveMessages
		stringArray := strings.Split(msg, "\x01")
		var userName = stringArray[2]
		var content = stringArray[3]
		var now = stringArray[4]
		fmt.Printf("%s %s: %s", now, userName, content)
		fmt.Println("")
		if strings.HasPrefix(msg, ":q") || strings.HasPrefix(msg, ":quit") {
			fmt.Printf("%s is leaving", userName)
		}
	}
}

func checkError(err error, funcName string) {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Fatal error:%s-----in func:%s", err.Error(), funcName)
		os.Exit(1)
	}
}
func main() {
	udpAddr, err := net.ResolveUDPAddr("udp4", "192.168.1.101:1200")
	checkError(err, "main")

	var c Client
	c.alive = true
	c.sendingMessageQueue = make(chan string)
	c.receiveMessages = make(chan string)

	u, err := uuid.NewV4()
	checkError(err, "main")
	c.userID = *u

	fmt.Println("input name: ")
	_, err = fmt.Scanln(&c.userName)
	checkError(err, "main")

	c.connection, err = net.DialUDP("udp", nil, udpAddr)
	checkError(err, "main")
	defer func() {
		_ = c.connection.Close()
	}()

	c.funcSendMessage("joined")

	go c.printMessage()
	go c.receiveMessage()

	go c.sendMessage()
	c.readInput()

	c.funcSendMessage("left")

	os.Exit(0)
}
