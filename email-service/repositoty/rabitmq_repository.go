package repositoty

import (
	"encoding/json"
	"fmt"
	"github.com/hasib-003/newsLetterMicroservice/email-service/model"
	"github.com/streadway/amqp"
	"log"
)

type RabbitMqRepository struct {
	Connection *amqp.Connection
}

func NewRabbitMqRepository(connection *amqp.Connection) *RabbitMqRepository {
	return &RabbitMqRepository{
		Connection: connection,
	}
}
func (r *RabbitMqRepository) ConsumeMessage() (<-chan model.UserWithNews, error) {
	if r.Connection == nil || r.Connection.IsClosed() {
		return nil, fmt.Errorf("RabbitMQ connection is not open")
	}
	log.Println("start consume message")
	channel, err := r.Connection.Channel()
	if err != nil {
		return nil, err
	}
	q, err := channel.QueueDeclare(
		"user_with_news",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err

	}
	msgs, err := channel.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	userWithNewsChan := make(chan model.UserWithNews)

	go func() {
		for d := range msgs {
			var usersNews model.UserWithNews
			err := json.Unmarshal(d.Body, &usersNews)
			if err != nil {
				log.Printf("unmarshal err:%v", err)
				continue
			}
			userWithNewsChan <- usersNews

			err = d.Ack(false)
			if err != nil {
				return
			}
		}
		close(userWithNewsChan)
	}()
	return userWithNewsChan, nil
}
