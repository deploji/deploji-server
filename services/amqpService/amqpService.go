package amqpService

import (
	"encoding/json"
	"github.com/sotomskir/mastermind-server/models"
	"github.com/streadway/amqp"
	"os"
)

func Send(deployment models.Deployment) error {
	url := os.Getenv("AMQP_URL")
	conn, err := amqp.Dial(url)
	if err != nil {
		return err
	}
	defer conn.Close()
	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	err = channel.Qos(1, 0, false)
	if err != nil {
		return err
	}
	defer channel.Close()
	queue, err := channel.QueueDeclare("deployments", true, false, false, false, nil)
	if err != nil {
		return err
	}
	body, err := json.Marshal(deployment)
	err = channel.Publish("", queue.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType: "text/plain",
		Body: body,
	})
	if err != nil {
		return err
	}
	return nil
}
