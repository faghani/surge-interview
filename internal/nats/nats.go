package nats

import (
	"github.com/nats-io/nats.go"
)

func New(adr string) (*nats.Conn, error) {
	nc, err := nats.Connect(adr)
	if err != nil {
		return nil, err
	}
	return nc, nil
}
