package netpoll

import (
	"context"
	"github.com/cloudwego/netpoll"
	"github.com/no-today/go-examples/netpoll/protocol"
	"log"
	"time"
)

func Server(network, address string) {
	listener, err := netpoll.CreateListener(network, address)
	if err != nil {
		log.Fatalf("Create listener [%s] %s failed, %v", network, address, err)
	}

	eventLoop, err := netpoll.NewEventLoop(
		handle,
		netpoll.WithOnPrepare(prepare),
		netpoll.WithReadTimeout(2*time.Second),
	)

	if err != nil {
		log.Fatalf("New eventLo [%s] %s failed, %v", network, address, err)
	}

	log.Printf("Start listen: [%s] %s\n", network, address)

	if err := eventLoop.Serve(listener); err != nil {
		log.Fatalf("Listen: [%s] %s failed, %v", network, address, err)
	}
}

func prepare(connection netpoll.Connection) context.Context {
	return context.Background()
}

func handle(ctx context.Context, connection netpoll.Connection) error {
	writer := connection.Writer()
	reader := connection.Reader()
	defer reader.Release()

	for {
		message, err := protocol.UnPackZeroCopy(reader)
		if err != nil {
			log.Printf("Read error, %v\n", err)
			return err
		}

		log.Printf("Received: %v\n", message)

		if "PING" == string(message.Content()) {
			bytes, _ := protocol.NewMessage(1, []byte("PONG")).Pack()
			_, err := writer.WriteBinary(bytes)
			if err != nil {
				log.Printf("Write error, %v\n", err)
				return err
			}

			writer.Flush()
		}
	}
}
