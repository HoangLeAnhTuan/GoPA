package broker

import (
	"fmt"
	"net"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func Open(url string) (*amqp091.Connection, error) {
	connection, err := amqp091.DialConfig(url, amqp091.Config{
		Dial: func(network, address string) (net.Conn, error) {
			return (&net.Dialer{Timeout: 3 * time.Second}).Dial(network, address)
		},
	})
	if err != nil {
		return nil, fmt.Errorf("connect rabbitmq: %w", err)
	}
	return connection, nil
}
