package Services

import (
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func HandleNotification() {

}

func StartConsumer(ws chan []byte) {
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Dialing AMQP Server")

	ch, err := conn.Channel()
	failOnError(err, "Creating AMQP Channel")
	defer func(ch *amqp091.Channel) {
		err := ch.Close()
		if err != nil {
			failOnError(err, "Closing Channel")
		}
	}(ch)

	err = ch.ExchangeDeclare(
		"checkrr", // name
		"fanout",  // type
		true,      // durable
		false,     // auto-deleted
		false,     // internal
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Declaring Exchange")

	q, err := ch.QueueDeclare("", false, false, true, false, nil)

	failOnError(err, "Declaring Queue")

	err = ch.QueueBind(
		q.Name,
		"",
		"checkrr",
		false,
		nil,
	)
	failOnError(err, "Binding Queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			HandleNotification()
			fmt.Println("Notification Receiver received: ", string(d.Body))
			ws <- d.Body
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}
