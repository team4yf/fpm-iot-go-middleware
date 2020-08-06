//Package tcp the tcp package
// for receive the tcp message, 
// only support the mcu2 env device
// it will remove after some other device is ok
// just a Read interface
package tcp

import (
	"net"
	"fmt"
	"bufio"
)


//ReceiveHandler the receive handler
//topic is the #sn/tcp/receive
//data  is the origin payload
type ReceiveHandler func(topic string, data []byte)

//NetReceiver the interface for the tcp receiver
//Read is the only one method should be implemented! it's just run the tcp server on a port
type NetReceiver interface {
	Read(port int)
}

type netReceiver struct {
	handler ReceiveHandler
}

//Read 
//start to open a socket server on the port
//use go routine to read, call the handler after received the message data.
func (r *netReceiver) Read(port int){
	go func(){
		// listen on all interfaces
		ln, _ := net.Listen("tcp", ":" + port)

		// accept connection on port
		conn, _ := ln.Accept()

		// run loop forever (or until ctrl-c)
		for {
			// will listen for message to process ending in newline (\n)
			message, _ := bufio.NewReader(conn).ReadString('\n')
			// output message received
			fmt.Print("Message Received:", string(message))
			r.handler("tes", message)
		}
	}()
	
}

func NewNetReceiver(f ReceiveHandler) NetReceiver{
	return &netReceiver{
		handler: f
	}
}