package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log_indexer/domain"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type Consumer struct {
	Reader    *kafka.Reader
	Dialer    *kafka.Dialer
	EsRepo    domain.EsInterface
	AlertRepo domain.AlertsInterface
	Topic     string
}

func (c *Consumer) CreateConnection() {
	godotenv.Load()
	var hostname = os.Getenv("HOST")
	var kafkaPort = os.Getenv("KAFKA_PORT")
	c.Reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{hostname + ":" + kafkaPort},
		Topic:     c.Topic,
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
		MaxWait:   time.Millisecond * 10,
		Dialer:    c.Dialer,
		GroupID:   "logs-consumer-group",
	})
	c.Reader.SetOffset(0)
}

func (c *Consumer) Start() {
	ctx := context.Background()
	for {
		msg, err := c.Reader.ReadMessage(ctx)
		if err != nil {
			logrus.Errorf("Failed to consume message from kafka topic %v", err)
			break
		}
		var readLog domain.Log
		err = json.Unmarshal(msg.Value, &readLog)
		if err != nil {
			logrus.Errorf("Failed to parse kafka message %v", err)
			break
		}
		c.checkForAlerts(ctx, &readLog, msg.Value)

		fmt.Printf("Read Message from Kafka: %+v\n", readLog)
		err = c.Reader.CommitMessages(ctx, msg)
		if err != nil {
			logrus.Errorf("Failed to commit kafka message %v", err)
			break
		}
		fmt.Printf("Kafka next for log %v", readLog.Message)
		c.EsRepo.IndexLog(ctx, readLog)
	}
}

func (c *Consumer) checkForAlerts(ctx context.Context, log *domain.Log, logBytes []byte) {
	if strings.ToLower(log.Level) == "error" {
		c.AlertRepo.TriggerAlert(ctx, logBytes)
	}
}
