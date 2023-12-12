package domain

import (
	"time"
)

type Log struct {
	Date    time.Time `json:"date"`
	Message string    `json:"message"`
	Level   string    `json:"level"`
}
