//Package tcp the tcp package
// for receive the tcp message,
// only support the mcu2 env device
// it will remove after some other device is ok
// just a Read interface
package tcp

import (
	"fmt"
	"net"
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
func (r *netReceiver) Read(port int) {
	go func() {
		// listen on all interfaces
		ln, _ := net.Listen("tcp", fmt.Sprintf(":%d", port))

		// run loop forever (or until ctrl-c)
		for {
			// accept connection on port
			conn, _ := ln.Accept()
			//logs an incoming message
			// Handle connections in a new goroutine.
			go func() {
				for {
					buf := make([]byte, 1024)
					reqLen, err := conn.Read(buf)
					if err != nil {
						fmt.Println("Error to read message because of ", err)
						return
					}
					protocol := buf[0]
					if protocol != 0xee {
						// only support protocol of startsWith ee
						return
					}
					data := buf[15:reqLen]
					// output message received
					r.handler("#socket/ee", data)
				}
			}()

		}
	}()

}

//NewNetReceiver create a new receiver
func NewNetReceiver(f ReceiveHandler) NetReceiver {
	return &netReceiver{
		handler: f,
	}
}
