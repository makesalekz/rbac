package data

import (
	"github.com/nats-io/nats.go"
	"gitlab.calendaria.team/services/rbac/internal/conf"
)

// NewNatsClient .
func NewNatsClient(conf *conf.Bootstrap) (*nats.Conn, func(), error) {
	nc, err := nats.Connect(conf.Nats)
	if err != nil {
		return nil, nil, err
	}

	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		nc.Close()
		return nil, nil, err
	}

	cleanup := func() {
		ec.Close()
		nc.Close()
	}

	return ec.Conn, cleanup, nil
}
