package main

import (
	"log_indexer/repository"
	"os"
	"time"

	es7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/joho/godotenv"
	kafka "github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

func main() {
	godotenv.Load()
	esConfig := es7.Config{
		Username:  "elastic",
		Password:  "changeme",
		Addresses: []string{os.Getenv("ELASTICSEARCH_HOST")},
	}
	esClient, err := es7.NewClient(esConfig)
	if err != nil {
		logrus.Errorf("Failed to open elasticsearch %v", err)
		time.Sleep(time.Second * 60)
	}
	esRepo := repository.NewEsRepo(esClient)

	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}
	logConsumer := repository.Consumer{
		Dialer: dialer,
		Topic:  os.Getenv("KAFKA_LOG_TOPIC_NAME"),
		EsRepo: &esRepo,
	}
	logConsumer.CreateConnection()
	logConsumer.Start()
}
