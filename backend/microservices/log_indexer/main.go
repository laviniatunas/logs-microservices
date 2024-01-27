package main

import (
	"log_indexer/repository"
	"os"
	"time"

	es7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	kafka "github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

func main() {
	defer func() {
		logrus.Info("Waiting 180 seconds before program exit")
		time.Sleep(time.Second * 180)
	}()
	godotenv.Load()
	esConfig := es7.Config{
		Username:  "elastic",
		Password:  "changeme",
		Addresses: []string{os.Getenv("ELASTICSEARCH_HOST")},
	}
	esClient, err := es7.NewClient(esConfig)
	if err != nil {
		logrus.Errorf("Failed to open elasticsearch %v", err)
		return
	}
	esRepo := repository.NewEsRepo(esClient)
	rmqConn, err := amqp.Dial(os.Getenv("RABBITMQ_HOST"))
	if err != nil {
		logrus.Errorf("Failed to open rabbitmq %v", err)
		return
	}
	alertRepo, err := repository.NewAlertsRepo(rmqConn)
	if err != nil {
		logrus.Errorf("Failed to open rabbitmq %v", err)
		return
	}
	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}
	logConsumer := repository.Consumer{
		Dialer:    dialer,
		Topic:     os.Getenv("KAFKA_LOG_TOPIC_NAME"),
		EsRepo:    &esRepo,
		AlertRepo: alertRepo,
	}
	logConsumer.CreateConnection()
	logConsumer.Start()
}
