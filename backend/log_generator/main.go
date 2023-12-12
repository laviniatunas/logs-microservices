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

func main() {
	godotenv.Load()
	var logFile = os.Getenv("LOG_FILE")
	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	possibleLogs := []string{"Detected malicious files: %v", "Program exited with error code %v", "Tried connectiong to host %v times", "Processed %v files", "Connection timed out after %v seconds"}

	for i := 0; i < 22; i++ {
		fmt.Printf("%v\n", i)
		code := rand.Intn(150)
		logMessage := possibleLogs[rand.Intn(len(possibleLogs))]
		logrus.SetOutput(f)
		logrus.Infof(logMessage, code)
		time.Sleep(time.Second * 3)
	}
}
