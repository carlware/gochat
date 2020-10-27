package rabbitmq

import (
	log "github.com/inconshreveable/log15"
	"github.com/streadway/amqp"
)

type server struct {
	host       string
	connection *amqp.Connection
}

func getExchangeName() string {
	return "chat"
}

func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
}

func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		getExchangeName(), // name
		"topic",           // type
		true,              // durable
		false,             // auto-deleted
		false,             // internal
		false,             // no-wait
		nil,               // arguments
	)
}

func NewServer(host string) (*server, error) {
	conn, err := amqp.Dial(host)

	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	defer ch.Close()

	err = declareExchange(ch)
	if err != nil {
		return nil, err
	}

	// Handle any errors if we were unable to create the queue
	if err != nil {
		return nil, err
	}

	return &server{
		host:       host,
		connection: conn,
	}, nil

}

func (p *server) Send(message []byte, topic string) error {
	channel, err := p.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	err = channel.Publish(
		getExchangeName(),
		topic,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		},
	)
	log.Info("Sending to RabbitMQ", "queue", topic, "data", string(message))
	return err
}

func (p *server) Listen(topic string) (chan []byte, error) {
	ch, err := p.connection.Channel()
	if err != nil {
		return nil, err
	}

	queue, err := declareRandomQueue(ch)
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(queue.Name, topic, getExchangeName(), false, nil)
	if err != nil {
		return nil, err
	}

	msgs, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	recv := make(chan []byte)
	go func() {
		for msg := range msgs {
			log.Info("Receiving to RabbitMQ", "queue", topic, "data", string(msg.Body))
			recv <- msg.Body
		}
	}()

	return recv, nil
}
