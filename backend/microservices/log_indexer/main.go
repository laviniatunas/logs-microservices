package main

import (
	"log_indexer/repository"
	"os"
	"time"

	"github.com/joho/godotenv"
	kafka "github.com/segmentio/kafka-go"
)

func main() {
	godotenv.Load()
	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}
	logConsumer := repository.Consumer{
		Dialer: dialer,
		Topic:  os.Getenv("KAFKA_LOG_TOPIC_NAME"),
	}
	logConsumer.CreateConnection()

	logConsumer.Start()
}
