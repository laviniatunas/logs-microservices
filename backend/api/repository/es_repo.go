package repository

import (
	"api/domain"
	"context"
	"encoding/json"
	"io"
	"reflect"
	"strings"

	es7 "github.com/elastic/go-elasticsearch/v7"
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

func (e *EsRepo) GetLogs(ctx context.Context) ([]domain.Log, error) {
	// query := `{"query": {"match_all" : {}}}`
	query := `{"query": { "match_all": {}}, "sort": [{"date": {"order": "desc"}}]}`
	var b strings.Builder
	b.WriteString(query)
	read := strings.NewReader(b.String())

	res, err := e.esClient.Search(
		e.esClient.Search.WithContext(ctx),
		e.esClient.Search.WithIndex("log_index"),
		e.esClient.Search.WithBody(read),
		e.esClient.Search.WithTrackTotalHits(true),
		e.esClient.Search.WithPretty(),
		e.esClient.Search.WithSize(500),
	)

	if err != nil {
		logrus.Errorf("Elasticsearch Search() API ERROR:", err)
		return nil, err
	} else {
		logrus.Println("res TYPE:", reflect.TypeOf(res))
	}

	buf, err := io.ReadAll(res.Body)
	if err != nil {
		logrus.Errorf("Error reading Body:", err)
		return nil, err
	}
	res.Body.Close()

	transformedResp := &domain.ElasticsearchResponse{}
	err = json.Unmarshal(buf, transformedResp)
	if err != nil {
		logrus.Errorf("Failed to unmarshal elasticsearch response %v", err)
		return nil, err
	}
	logs := make([]domain.Log, 0)
	for _, log := range transformedResp.Hits.Hits {
		logs = append(logs, log.Source)
	}
	logrus.Printf("String is %+v", logs)

	return logs, nil

	// logBytes, err := json.Marshal(log)
	// if err != nil {
	// 	logrus.Errorf("Failed to marshal log message %v", err)
	// 	return fmt.Errorf("Failed to marshal log message %v", err)
	// }
	// logrus.Infof("Logging %v", log.Message)
	// indexRequest := esapi.IndexRequest{
	// 	Index:      EsIndex,
	// 	DocumentID: strconv.FormatUint(uint64(time.Now().Unix()), 10),
	// 	Body:       bytes.NewReader(logBytes),
	// 	Refresh:    "true",
	// }
	// response, err := indexRequest.Do(ctx, e.esClient)
	// if err != nil {
	// 	logrus.Errorf("Failed to do index request %v", err)
	// 	return fmt.Errorf("Failed to do index request %v", err)
	// }
	// logrus.Infof("Response is = %v", response.StatusCode)
}
