package rabbitmq

import (
	"fmt"

	"github.com/streadway/amqp"
)

type server struct {
	queueName string
	host      string
	channel   *amqp.Channel
	queue     *amqp.Queue
}

func NewServer(host, queueName string) (*server, error) {
	conn, err := amqp.Dial(host)

	if err != nil {
		fmt.Println(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	defer ch.Close()

	queue, err := ch.QueueDeclare(queueName, false, false, false, false, nil)

	// Handle any errors if we were unable to create the queue
	if err != nil {
		return nil, err
	}

	return &server{
		queueName: queueName,
		host:      host,
		channel:   ch,
		queue:     &queue,
	}, nil

}

func (p *server) Sender(message []byte) error {
	err := p.channel.Publish(
		"",
		p.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		},
	)
	return err
}

func (p *server) Listen() (chan []byte, error) {
	recv := make(chan []byte)
	msgs, err := p.channel.Consume(p.queueName, "", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	go func() {
		for msg := range msgs {
			recv <- msg.Body
		}
	}()
	return recv, nil
}
