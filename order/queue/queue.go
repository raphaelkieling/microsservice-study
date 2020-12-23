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

func StartConsuming(ch *amqp.Channel, in chan []byte) {
	q, err := ch.QueueDeclare(
		"checkout_queue",
		true,
		false,
		false,
		false,
		nil,
	)

	fmt.Println("Declarando queue")

	if err != nil {
		panic(err.Error())
	}

	msgs, err := ch.Consume(q.Name, "checkout", true, false, false, false, nil)

	fmt.Printf("Consumindo %v mensagens", len(msgs))

	if err != nil {
		panic(err.Error())
	}

	go func() {
		for m := range msgs {
			fmt.Println("Resgatando dados")
			in <- []byte(m.Body)
		}
		close(in)
	}()
}
