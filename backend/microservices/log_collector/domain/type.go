package domain

import (
	"time"

	kafka "github.com/segmentio/kafka-go"
)

type Log struct {
	Date    time.Time `json:"date"`
	Message string    `json:"message"`
	Level   string    `json:"level"`
}

type Producer struct {
	Writer *kafka.Writer
	Dialer *kafka.Dialer
}
