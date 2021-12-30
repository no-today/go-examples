package native

import (
	"github.com/no-today/go-examples/netpoll/protocol"
	"io"
	"log"
	"net"
)

func Server(network, address string) {
	listen, err := net.Listen(network, address)
	if err != nil {
		log.Fatalf("Listen: [%s] %s failed, %v", network, address, err)
	}
	defer listen.Close()

	log.Printf("Start listen: [%s] %s\n", network, address)

	for {
		// 建立连接都是使用的同一个 goroutine
		conn, err := listen.Accept()
		if err != nil {
			continue
		}

		// 为每个请求分配一个 goroutine 去处理业务逻辑
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

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

		if string(message.Content()) == "PING" {
			bytes, _ := protocol.NewMessage(1, []byte("PONG")).Pack()
			_, err = conn.Write(bytes)
			if err != nil {
				log.Printf("Write error, %v\n", err)
				break
			}
		}
	}
}
