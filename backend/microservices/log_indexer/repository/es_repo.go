package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log_indexer/domain"
	"strconv"
	"time"

	es7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/sirupsen/logrus"
)

const (
	EsIndex = "log_index"
)

type EsRepo struct {
	esClient *es7.Client
}

func NewEsRepo(esClient *es7.Client) EsRepo {
	return EsRepo{
		esClient: esClient,
	}
}

func (e *EsRepo) IndexLog(ctx context.Context, log domain.Log) error {
	logBytes, err := json.Marshal(log)
	if err != nil {
		logrus.Errorf("Failed to marshal log message %v", err)
		return fmt.Errorf("Failed to marshal log message %v", err)
	}
	logrus.Infof("Logging %v", log.Message)
	indexRequest := esapi.IndexRequest{
		Index:      EsIndex,
		DocumentID: strconv.FormatUint(uint64(time.Now().Unix()), 10),
		Body:       bytes.NewReader(logBytes),
		Refresh:    "true",
	}
	response, err := indexRequest.Do(ctx, e.esClient)
	if err != nil {
		logrus.Errorf("Failed to do index request %v", err)
		return fmt.Errorf("Failed to do index request %v", err)
	}
	logrus.Infof("Response is = %v", response.StatusCode)
	return nil
}
