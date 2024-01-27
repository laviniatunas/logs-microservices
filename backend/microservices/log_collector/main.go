package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log_collector/domain"
	"net"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/hpcloud/tail"
	"github.com/joho/godotenv"
	kafka "github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

func main() {
	defer func() {
		logrus.Infof("Waiting 180 before program exit")
		time.Sleep(time.Second * 180)
	}()
	godotenv.Load()
	var logFile = os.Getenv("LOG_FILE")
	logTopic := os.Getenv("KAFKA_LOG_TOPIC_NAME")

	t, err := tail.TailFile(logFile, tail.Config{Follow: true})
	if err != nil {
		logrus.Fatalf("Could not start program, error: %v", err)
	}

	connection, err := kafka.Dial("tcp", net.JoinHostPort(os.Getenv("HOST"), os.Getenv("KAFKA_PORT")))
	if err != nil {
		logrus.Fatalf("Could not connect to Kafka, error: %v", err.Error())
	}

	logTopicConfig := kafka.TopicConfig{Topic: logTopic, NumPartitions: 1, ReplicationFactor: 1}
	err = connection.CreateTopics(logTopicConfig)
	if err != nil {
		logrus.Fatalf("Could not create Kafka topic, error: %v", err.Error())
	}

	createTopic(logTopic)
	producer := NewProducer()

	for line := range t.Lines {
		jsonBytes, err := generateJson(line.Text)
		if err != nil {
			continue
		}
		fmt.Println(string(jsonBytes))
		Produce(logTopic, producer, jsonBytes)
	}

}

func generateJson(logLine string) ([]byte, error) {
	exprTime, err := regexp.Compile("time=\"(\\w|[-:+])*\"")
	if err != nil {
		logrus.Errorf("Could not compile time regex, error: %v", err)
		return nil, err
	}

	timeStr := exprTime.FindString(logLine)
	parsedTime, err := time.Parse("2006-01-02T15:04:05", strings.Split(strings.Split(timeStr, "\"")[1], "Z")[0])
	if err != nil {
		logrus.Errorf("Could not find time regex, error: %v", err)
		return nil, err
	}

	exprLog, err := regexp.Compile("level=([^\\s]+)")
	if err != nil {
		logrus.Errorf("Could not compile log level regex, error: %v", err)
		return nil, err
	}

	logStr := exprLog.FindString(logLine)
	parsedLog := strings.Split(logStr, "=")[1]

	exprMsg, err := regexp.Compile("msg=\".*\"")
	if err != nil {
		logrus.Errorf("Could not compile message regex, error: %v", err)
		return nil, err
	}

	msgStr := exprMsg.FindString(logLine)
	parsedMsg := strings.Split(msgStr, "\"")[1]

	lineLog := domain.Log{
		Date:    parsedTime,
		Message: parsedMsg,
		Level:   parsedLog,
	}
	jsonBytes, _ := json.Marshal(lineLog)
	return jsonBytes, nil
}

func createTopic(logTopic string) {
	connection, err := kafka.Dial("tcp", net.JoinHostPort(os.Getenv("HOST"), os.Getenv("KAFKA_PORT")))
	if err != nil {
		logrus.Fatalf("Could not connect to Kafka, error: %v", err.Error())
	}

	logTopicConfig := kafka.TopicConfig{Topic: logTopic, NumPartitions: 1, ReplicationFactor: 1}
	err = connection.CreateTopics(logTopicConfig)
	if err != nil {
		logrus.Fatalf("Could not create Kafka topic, error: %v", err.Error())
	}
}

func Produce(topic string, producer *domain.Producer, value []byte) {
	err := producer.Writer.WriteMessages(context.Background(), kafka.Message{
		Topic:  topic,
		Offset: 0,
		Value:  value,
	})

	if err != nil {
		fmt.Printf("delivery failed %s \n", err.Error())
	} else {
		fmt.Printf("message delivered topic: %s \n", topic)
	}
}

func NewProducer() *domain.Producer {
	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{os.Getenv("HOST") + ":" + os.Getenv("KAFKA_PORT")},
		Dialer:  dialer,
	})

	return &domain.Producer{
		Writer: writer,
	}
}
