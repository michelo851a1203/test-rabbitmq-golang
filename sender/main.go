package main

import (
	"log"

	"github.com/streadway/amqp"
)

func main() {
	amqpConnection, err := amqp.Dial("amqp://user:secret@localhost:5672")
	if err != nil {
		log.Println("connection error")
		panic(err)
	}
	defer func() {
		amqpConnection.Close()
	}()

	// build channel
	amqpChannel, err := amqpConnection.Channel()
	if err != nil {
		log.Println("advanced message queue protocol error : create channel error")
		panic(err)
	}

	// build queue with rabbit mq on golang application
	mainQueue, err := amqpChannel.QueueDeclare(
		"main_queue",
		false, // durable
		false, // delete when used
		false, // exclusive
		false, // no wait
		nil,
	)
	if err != nil {
		log.Println("build queue error")
		panic(err)
	}
	// here try to transfer data
	body := "hello world"
	err = amqpChannel.Publish(
		"", // exchange
		mainQueue.Name,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	if err != nil {
		log.Println("publish error")
		panic(err)
	}

}
