//Package tcp the tcp package
// for receive the tcp message,
// only support the mcu2 env device
// it will remove after some other device is ok
// just a Read interface
package tcp

import (
	"fmt"
	"net"
	"github.com/team4yf/fpm-iot-go-middleware/pkg/utils"
)

//ReceiveHandler the receive handler
//clientID is the random id of the tcp connection
//data  is the origin payload
type ReceiveHandler func(clientID string, data []byte)

//NetReceiver the interface for the tcp receiver
//Read is the only one method should be implemented! it's just run the tcp server on a port
type NetReceiver interface {
	Read(max int)
	Write(clientID string, buf []byte) error
	Listen(port int)
}

type netReceiver struct {
	handler ReceiveHandler
	clients map[string]net.Conn
}

//Read
//start to open a socket server on the port
//use go routine to read, call the handler after received the message data.
func (r *netReceiver) Listen(port int) {
	go func() {
		// listen on all interfaces
		ln, _ := net.Listen("tcp", fmt.Sprintf(":%d", port))

		// run loop forever (or until ctrl-c)
		for {
			// accept connection on port
			conn, _ := ln.Accept()
			clientID := utils.GenShortID()
			r.clients[clientID] = conn
		}
	}()

}

func (r *netReceiver) Read(max int) {
		
	// run loop forever (or until ctrl-c)
	for clientID, conn := range r.clients {
		//logs an incoming message
		// Handle connections in a new goroutine.
		go func() {
			for {
				buf := make([]byte, max)
				reqLen, err := conn.Read(buf)
				if err != nil {
					fmt.Println("Error to read message because of ", err)
					return
				}
				
				data := buf[15:reqLen]
				// output message received
				r.handler(clientID, data)
			}
		}()

	}

}


func (r *netReceiver) Write(clientID string, buf []byte) error {
	return nil
}


//NewNetReceiver create a new receiver
func NewNetReceiver(f ReceiveHandler) NetReceiver {
	return &netReceiver{
		handler: f,
	}
}
