package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/streadway/amqp"
)

func main() {
	amqpConnection, err := amqp.Dial("amqp://user:secret@localhost:5672")
	defer func() {
		amqpConnection.Close()
	}()
	if err != nil {
		log.Println("connection error")
		panic(err)
	}

	amqpChannel, err := amqpConnection.Channel()

	defer func() {
		amqpChannel.Close()
	}()

	if err != nil {
		log.Println("amqp create channel error")
		panic(err)
	}
	mainQueue, err := amqpChannel.QueueDeclare(
		"main_queue",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Println("main queue error")
		panic(err)
	}
	// get data from comsumer
	mainMessageChannel, err := amqpChannel.Consume(
		mainQueue.Name, // queue name
		"",             // consumer
		true,           // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,
	)

	if err != nil {
		log.Println("consume message error")
		panic(err)
	}

	go func() {
		for messageInfo := range mainMessageChannel {
			log.Printf("received message : %s", messageInfo.Body)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	<-quit
}
