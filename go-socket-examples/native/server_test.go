package native

import (
	"testing"
)

const (
	network = "tcp"
	address = ":8090"
)

func TestServer(t *testing.T) {
	Server(network, address)
}

func TestClient(t *testing.T) {
	Client(network, address)
}
