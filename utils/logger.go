package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	LogWarning *log.Logger
	LogInfo    *log.Logger
	LogError   *log.Logger
	LogDebug   *log.Logger
	LogDefault *log.Logger
)

func ToString(i interface{}) string {
	log, _ := json.Marshal(i)
	logString := string(log)

	return logString
}

func EndLogPrint(str interface{}) {
	log.Println(fmt.Sprintf("%s %s %s", strings.Repeat("-", 15), str, strings.Repeat("-", 15)))
}

func InitLogger() {
	file, err := os.OpenFile(fmt.Sprintf("%s/%s", "logs", "gologs.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	LogWarning = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	LogInfo = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	LogError = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	LogDebug = log.New(file, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	LogDefault = log.New(file, "LOG", log.Ldate|log.Ltime|log.Lshortfile)
}
