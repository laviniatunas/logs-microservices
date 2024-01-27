package main

import (
	"api/repository"
	"context"
	"io"
	"net/http"
	"os"
	"time"

	es7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

var db = make(map[string]string)

var ClientChan = make(chan string)

func setupRouter(esRepo repository.EsRepo) *gin.Engine {
	r := gin.Default()

	r.GET("/logs", enableCORS(), func(c *gin.Context) {
		logs, err := esRepo.GetLogs(c)
		if err != nil {
			logrus.Errorf("Error in /logs call %v", err)
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusOK,
				"message": logs,
			})
		}

	})

	r.GET("/stream", HeadersMiddleware(), enableCORS(), func(c *gin.Context) {
		isClientDced := c.Stream(func(w io.Writer) bool {
			// Stream message to client from message channel
			logrus.Infof("Waiting receive on client chan %v", ClientChan)
			if msg, ok := <-ClientChan; ok {
				c.SSEvent("message", msg)
				return true
			}
			return false
		})
		if isClientDced {
			logrus.Errorf("Client disconnected")
		}
	})

	return r
}

func main() {
	defer func() {
		logrus.Info("Waiting for 180 seconds")
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
		time.Sleep(time.Second * 60)
	}
	esRepo := repository.NewEsRepo(esClient)

	rmqConn, err := amqp.Dial(os.Getenv("RABBITMQ_HOST"))
	if err != nil {
		logrus.Errorf("Failed to open rabbitmq %v", err)
		return
	}

	go startAlertsRepo(rmqConn)

	r := setupRouter(esRepo)
	// // Listen and Server in 0.0.0.0:8080
	r.Run(":8000")
}

func startAlertsRepo(rmqConn *amqp.Connection) {
	alertsRepo, err := repository.NewAlertsRepo(rmqConn, ClientChan)
	if err != nil {
		logrus.Errorf("Failed to open rabbit mq %v", err)
		return
	}
	alertsRepo.Start(context.Background())
}

func HeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")
		c.Next()
	}
}
func enableCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		c.Next()
	}
}
