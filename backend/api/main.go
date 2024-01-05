package main

import (
	"api/repository"
	"net/http"
	"os"
	"time"

	es7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var db = make(map[string]string)

func setupRouter(esRepo repository.EsRepo) *gin.Engine {
	r := gin.Default()

	r.GET("/logs", func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

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

	return r
}

// func setupRouter() *gin.Engine {
// 	// Disable Console Color
// 	// gin.DisableConsoleColor()
// 	r := gin.Default()

// 	// Ping test
// 	r.GET("/ping", func(c *gin.Context) {
// 		c.String(http.StatusOK, "pong")
// 	})

// 	// Get user value
// 	r.GET("/user/:name", func(c *gin.Context) {
// 		user := c.Params.ByName("name")
// 		value, ok := db[user]
// 		if ok {
// 			c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
// 		} else {
// 			c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
// 		}
// 	})

// 	// Authorized group (uses gin.BasicAuth() middleware)
// 	// Same than:
// 	// authorized := r.Group("/")
// 	// authorized.Use(gin.BasicAuth(gin.Credentials{
// 	//	  "foo":  "bar",
// 	//	  "manu": "123",
// 	//}))
// 	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
// 		"foo":  "bar", // user:foo password:bar
// 		"manu": "123", // user:manu password:123
// 	}))

// 	/* example curl for /admin with basicauth header
// 	   Zm9vOmJhcg== is base64("foo:bar")

// 		curl -X POST \
// 	  	http://localhost:8080/admin \
// 	  	-H 'authorization: Basic Zm9vOmJhcg==' \
// 	  	-H 'content-type: application/json' \
// 	  	-d '{"value":"bar"}'
// 	*/
// 	authorized.POST("admin", func(c *gin.Context) {
// 		user := c.MustGet(gin.AuthUserKey).(string)

// 		// Parse JSON
// 		var json struct {
// 			Value string `json:"value" binding:"required"`
// 		}

// 		if c.Bind(&json) == nil {
// 			db[user] = json.Value
// 			c.JSON(http.StatusOK, gin.H{"status": "ok"})
// 		}
// 	})

// 	return r
// }

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

	r := setupRouter(esRepo)
	// // Listen and Server in 0.0.0.0:8080
	r.Run(":8000")
}
