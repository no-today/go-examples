package netpoll

import (
	"github.com/cloudwego/netpoll"
	"github.com/no-today/go-examples/netpoll/protocol"
	"io"
	"log"
	"time"
)

func Client(network, address string, timeout time.Duration) {
	conn, err := netpoll.DialConnection(network, address, timeout)
	if err != nil {
		log.Fatalf("Connection: [%s] %s failed, %v", network, address, err)
	}

	log.Printf("Connected [%s] %s\n", network, address)

	go ping(conn, 1*time.Second)

	reader := conn.Reader()
	for {
		message, err := protocol.UnPackZeroCopy(reader)
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Printf("Read error, %v\n", err)
			return
		}

		log.Printf("Received: %v\n", message)
	}
}

func ping(conn netpoll.Connection, duration time.Duration) {
	writer := conn.Writer()

	for {
		bytes, _ := protocol.NewMessage(1, []byte("PING")).Pack()

		// 强行把一个包分开发
		_, err := writer.WriteBinary(bytes[:3])
		if err != nil {
			log.Println("Write error")
			return
		}

		time.Sleep(duration)

		_, err = writer.WriteBinary(bytes[3:])
		if err != nil {
			log.Println("Write error")
			return
		}

		writer.Flush()
	}
}
