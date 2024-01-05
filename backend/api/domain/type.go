package domain

import (
	"time"
)

type Log struct {
	Date    time.Time `json:"date"`
	Message string    `json:"message"`
	Level   string    `json:"level"`
}

type ElasticsearchResponse struct {
	Hits struct {
		Hits []struct {
			Source Log `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}
