package broker

import (
	"encoding/json"
	"github.com/rabbitmq/amqp091-go"
	"log/slog"
	"os"
)

type RabbitBroker struct {
	channel *amqp091.Channel
	conn    *amqp091.Connection
}

func failOnError(err error, msg string) {
	if err != nil {
		slog.Error("RabbitBroker error: ", "error", err, "msg", msg)
		os.Exit(1)
	}
}

func NewRabbitBroker() *RabbitBroker {
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	err = ch.ExchangeDeclare(
		"like_exchange", // имя обмена
		"direct",        // тип обмена (можно использовать "fanout" для широковещательной рассылки)
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // аргументы
	)
	failOnError(err, "Failed to declare an exchange")

	_, err = ch.QueueDeclare(
		"new_like", // имя очереди
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		"new_like",      // имя очереди
		"new_like",      // routing key
		"like_exchange", // имя обмена
		false,
		nil,
	)
	failOnError(err, "Failed to bind queue to exchange")

	return &RabbitBroker{conn: conn, channel: ch}
}

func (b *RabbitBroker) PublishNewLike(postID, likerID, likedID int) error {

	body := struct {
		PostID int `json:"post_id"`
		Liker  int `json:"liker"`
		Liked  int `json:"liked"`
	}{PostID: postID, Liker: likerID, Liked: likedID}

	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return err
	}
	err = b.channel.Publish(
		"like_exchange",
		"new_like",
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        bodyJSON,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (b *RabbitBroker) Close() {
	if b.channel != nil {
		b.channel.Close()
	}
	if b.conn != nil {
		b.conn.Close()
	}
}
