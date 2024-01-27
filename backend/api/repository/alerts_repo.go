package repository

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type AlertsRepo struct {
	conn       *amqp.Connection
	rmqChannel *amqp.Channel
	rmqQueue   amqp.Queue
	ClientChan chan string
}

func NewAlertsRepo(conn *amqp.Connection, ClientChan chan string) (*AlertsRepo, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	rmqQueue, err := channel.QueueDeclare(
		"alerts_queue",
		true, false, false, false, nil,
	)
	return &AlertsRepo{
		conn:       conn,
		rmqChannel: channel,
		rmqQueue:   rmqQueue,
		ClientChan: ClientChan,
	}, nil
}

func (a *AlertsRepo) Start(ctx context.Context) error {
	msgs, err := a.rmqChannel.Consume(
		"alerts_queue", // queue
		"",             // consumer
		true,           // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	if err != nil {
		return err
	}
	for d := range msgs {
		logrus.Infof("Received a message: %s", d.Body)
		go func() {
			logrus.Infof("Trying to send to %v", a.ClientChan)
			a.ClientChan <- string(d.Body)
			logrus.Infof("Sent to another channel")
		}()
	}
	return nil
}
