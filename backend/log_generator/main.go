package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type LogMessage struct {
	Message  string
	LogLevel string
}

var (
	possibleLogs = []LogMessage{
		{
			Message:  "Detected malicious files: %v",
			LogLevel: "INFO",
		},
		{
			Message:  "Program exited with error code %v",
			LogLevel: "ERROR",
		},
		{
			Message:  "Tried connectiong to host %v times",
			LogLevel: "INFO",
		},
		{
			Message:  "Processed %v files",
			LogLevel: "INFO",
		},
		{
			Message:  "Connection timed out after %v seconds",
			LogLevel: "ERROR",
		},
		{
			Message:  "Indexer failed with error %v",
			LogLevel: "ERROR",
		},
		{
			Message:  "Something went wrong, please retry in %v seconds",
			LogLevel: "INFO",
		},
		{
			Message:  "Too many requests at the same time. Please wait %v seconds before trying another request",
			LogLevel: "INFO",
		},
	}
)

func main() {
	godotenv.Load()
	var logFile = os.Getenv("LOG_FILE")
	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	for {
		fmt.Printf("New log available")
		code := rand.Intn(150)
		logMessage := possibleLogs[rand.Intn(len(possibleLogs))]
		logrus.SetOutput(f)
		if logMessage.LogLevel == "ERROR" {
			logrus.Errorf(logMessage.Message, code)
		} else {
			logrus.Infof(logMessage.Message, code)
		}
		time.Sleep(time.Second * 3)
	}
}
