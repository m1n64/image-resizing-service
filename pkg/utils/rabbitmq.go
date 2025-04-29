package utils

import (
	"github.com/streadway/amqp"
	"log"
	"sync"
)

// deprecated
// Don't use

type RabbitMQConnection struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

var instance *RabbitMQConnection
var once sync.Once

func ConnectRabbitMQ(rabbitHost string) *RabbitMQConnection {
	once.Do(func() {
		conn, err := amqp.Dial(rabbitHost)
		if err != nil {
			log.Fatalf("Failed to connect to RabbitMQ: %s", err)
		}

		ch, err := conn.Channel()
		if err != nil {
			log.Fatalf("Failed to open a channel: %s", err)
		}

		instance = &RabbitMQConnection{
			Connection: conn,
			Channel:    ch,
		}
		log.Println("RabbitMQ connection established")
	})
	return instance
}

func InitializeQueues() {
	queues := []string{}

	for _, queue := range queues {
		_, err := instance.Channel.QueueDeclare(
			queue,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Fatalf("Failed to declare queue %s: %s", queue, err)
		}
		log.Printf("Queue declared: %s", queue)
	}
}

func GetRabbitMQInstance() *RabbitMQConnection {
	if instance == nil {
		log.Fatalf("RabbitMQ connection is not initialized. Call ConnectRabbitMQ first.")
	}
	return instance
}

func (r *RabbitMQConnection) CloseRabbitMQ() {
	if r.Channel != nil {
		r.Channel.Close()
	}
	if r.Connection != nil {
		r.Connection.Close()
	}
	log.Println("RabbitMQ connection closed")
}
