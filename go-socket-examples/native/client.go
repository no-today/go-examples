package native

import (
	"github.com/no-today/go-examples/netpoll/protocol"
	"io"
	"log"
	"net"
	"time"
)

func Client(network, address string) {
	conn, err := net.DialTimeout(network, address, time.Second)
	if err != nil {
		time.Sleep(3 * time.Second)
		Client(network, address)
		return
	}
	defer conn.Close()

	log.Printf("Connected [%s] %s\n", network, address)

	go ping(conn, 1*time.Second)

	for {
		// 阻塞到有一个完整的包才返回
		message, err := protocol.UnPack(conn)
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

func ping(conn net.Conn, duration time.Duration) {
	for {
		bytes, _ := protocol.NewMessage(1, []byte("PING")).Pack()

		// 强行把一个包分开发
		_, err := conn.Write(bytes[:2])
		if err != nil {
			log.Printf("Write error, %v\n", err)
			break
		}

		time.Sleep(duration)

		_, err = conn.Write(bytes[2:])
		if err != nil {
			log.Printf("Write error, %v\n", err)
			break
		}
	}
}
