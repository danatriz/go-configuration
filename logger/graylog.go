package logger

import (
	"configurations/dtos"
	"fmt"
	"os"
	"time"

	"github.com/Devatoria/go-graylog"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	LEVEL_ERROR = 4
	LEVEL_WARN  = 5
	LEVEL_INFO  = 6
	LEVEL_DEBUG = 7
)

const (
	NAME_ERROR = "ERROR"
	NAME_WARN  = "WARN"
	NAME_INFO  = "INFO"
	NAME_DEBUG = "DEBUG"
)

const (
	RETRY = 5
	SIZE  = 6400000
)

var GrayLogger *graylog.Graylog


func InitalzeGraylog() error {
	address := viper.GetString("graylog.address")
	graylog, err := graylog.NewGraylog(graylog.Endpoint{
		Transport: graylog.UDP,
		Address: address,
		Port: uint(viper.GetInt("graylog.port")),
	})

	if err != nil {
		return err
	}

	GrayLogger = graylog
	return nil
}

func GraylogRequest(shortMessage string, response string, context *gin.Context) {
	var headers, body string

	headerContext := context.Value("headerRequest")
	if headerContext != nil {
		headers = headerContext.(string)
	}

	bodyContext := context.Value("bodyRequest")
	if bodyContext != nil {
		body = bodyContext.(string)
	}

	fullMessage := fmt.Sprintf("Headers: %s\nBody: %s\nResponse: %s", headers, body, response)

	host, err := os.Hostname()
	log.Debug(host)
	if err != nil {
		os.Exit(1)
	}

	version := viper.GetString("configuration.version")
	environtment := viper.GetString("configuration.environment")
	serviceName := viper.GetString("service_name")
	fmt.Println(version, environtment, serviceName)

	if len(fullMessage) > SIZE {
		fullMessage = fullMessage[:SIZE]
	}
}

func GraylogStatus(request *dtos.GraylogMessage) {
	var levelInt int
	var levelName string

	switch request.Level {
	case "ERROR":
		levelInt = LEVEL_ERROR
		levelName = NAME_ERROR
	case "WARN":
		levelInt = LEVEL_WARN
		levelName = NAME_WARN
	case "DEBUG":
		levelInt = LEVEL_DEBUG
		levelName = NAME_DEBUG
	default:
		levelInt = LEVEL_INFO
		levelName = NAME_INFO
	}

	graylogMessage := graylog.Message{
		Version: request.Version,
		Host: request.Host,
		ShortMessage: request.ShortMessage,
		FullMessage: request.FullMessage,
		Timestamp: time.Now().Unix(),
		Level: uint(levelInt),
		Extra: map[string]string{
			"level_name":   levelName,
			"env":          request.Environment,
			"service_name": request.ServiceName,
		},
	}

	for i := 1; i < RETRY; i++ {
		err := GrayLogger.Send(graylogMessage)
		if err == nil {
			break
		}
	}
}