package broker

import (
	"github.com/rabbitmq/amqp091-go"
	"log"
	"log/slog"
	"pictureloader/notification_microservice/notifications/likes"
)

type RabbitMQ struct {
	Connection   *amqp091.Connection
	Channel      *amqp091.Channel
	notifService *likes.NotificationService
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func NewRabbitMQ(notifService *likes.NotificationService) *RabbitMQ {
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")

	// Открываем канал
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	return &RabbitMQ{conn, ch, notifService}
}

func (rmq *RabbitMQ) ListenLikes() {
	// Декларируем очередь
	q, err := rmq.Channel.QueueDeclare(
		"new_like", // имя очереди
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = rmq.Channel.QueueBind(
		"new_like",      // имя очереди
		"new_like",      // routing key
		"like_exchange", // имя обмена
		false,
		nil,
	)
	failOnError(err, "Failed to bind queue to exchange")

	// Регистрируем потребителя
	msgs, err := rmq.Channel.Consume(
		q.Name, // имя очереди
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	// Чтение сообщений из очереди
	for d := range msgs {
		slog.Info("Received a like message", "message body", d.Body)
		err := rmq.notifService.ProcessLikeMessage(d.Body)
		if err != nil {
			slog.Error("Failed to process like messages", "error", err)
		}
	}
}
