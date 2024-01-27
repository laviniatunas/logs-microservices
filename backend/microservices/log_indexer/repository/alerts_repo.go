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
}

func NewAlertsRepo(conn *amqp.Connection) (*AlertsRepo, error) {
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
	}, nil
}

func (a *AlertsRepo) TriggerAlert(ctx context.Context, logBytes []byte) error {
	message := amqp.Publishing{
		ContentType: "application/json",
		Body:        logBytes,
	}
	err := a.rmqChannel.PublishWithContext(ctx, "", a.rmqQueue.Name, false, false, message)
	if err != nil {
		logrus.Errorf("Failed to send an error message to rabbitMQ %v", err)
	}
	logrus.Infof("Sent an error message to rabbitMQ successfully")
	return err
}
