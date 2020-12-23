package queue

import (
	"fmt"

	"github.com/streadway/amqp"
)

func Connect() *amqp.Channel {
	dsn := "amqp://rabbitmq:rabbitmq@localhost:5672/"
	conn, err := amqp.Dial(dsn)

	if err != nil {
		panic(err.Error())
	}

	channel, err := conn.Channel()
	if err != nil {
		panic(err.Error())
	}

	return channel
}

func Notify(payload []byte, exchange string, routingKey string, ch *amqp.Channel) {
	err := ch.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(payload),
		})

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Sended")
}
