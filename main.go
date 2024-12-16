package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	fmt.Println("==============================")
	fmt.Println("project: emqx-rabbit")
	fmt.Println("usecase: consumer")
	fmt.Println("release: 12.12.24")
	fmt.Println("==============================")

	conn, err := amqp.Dial("amqp://admin:d@mn@localhost:5672/")
	failOnError("Failed to connect to RabbitMQ", err)
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError("Failed to open a channel", err)
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"trafficlights", // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	failOnError("Failed to declare a queue", err)

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError("Failed to register a consumer", err)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			fmt.Println()
			fmt.Printf("RoutingKey:%s\n", d.RoutingKey)
			fmt.Printf("Body      :%s\n", d.Body)
		}
	}()

	fmt.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func failOnError(msg string, err error) {
	if err != nil {
		log.Fatalf("[%s]: %s", msg, err.Error())
	}
}
