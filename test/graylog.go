package test

import (
	"configurations/dtos"
	"configurations/logger"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

func TestGraylog() (string, error) {
	// Setup configuration
	viper.SetDefault("graylog.address", "localhost")
	viper.SetDefault("graylog.port", 12201)
	viper.SetDefault("configuration.version", "1.0.0")
	viper.SetDefault("configuration.environment", "development")
	viper.SetDefault("service_name", "exam-service")

	// Initialize Graylog
	err := logger.InitalzeGraylog()
	if err != nil {
		return "", err
	}

	log.Println("Graylog initialized successfully")

	// Get hostname
	host, err := os.Hostname()
	if err != nil {
		log.Printf("Failed to get hostname: %v", err)
		host = "localhost"
	}

	// Test sending different log levels
	testMessages := []dtos.GraylogMessage{
		{
			Version:      "1.0.0",
			Host:         host,
			ShortMessage: "Test INFO message",
			FullMessage:  "This is a test INFO message from Go application",
			Environment:  "development",
			ServiceName:  "exam-service",
			Level:        "INFO",
		},
		{
			Version:      "1.0.0",
			Host:         host,
			ShortMessage: "Test WARN message",
			FullMessage:  "This is a test WARNING message from Go application",
			Environment:  "development",
			ServiceName:  "exam-service",
			Level:        "WARN",
		},
		{
			Version:      "1.0.0",
			Host:         host,
			ShortMessage: "Test ERROR message",
			FullMessage:  "This is a test ERROR message from Go application",
			Environment:  "development",
			ServiceName:  "exam-service",
			Level:        "ERROR",
		},
		{
			Version:      "1.0.0",
			Host:         host,
			ShortMessage: "Test DEBUG message",
			FullMessage:  "This is a test DEBUG message from Go application",
			Environment:  "development",
			ServiceName:  "exam-service",
			Level:        "DEBUG",
		},
	}

	// Send test messages
	for i, msg := range testMessages {
		log.Printf("Sending test message %d: %s", i+1, msg.ShortMessage)
		logger.GraylogStatus(&msg)
		time.Sleep(500 * time.Millisecond) // Wait between messages
	}

	result := "All test messages sent to Graylog!\nCheck Graylog web interface at http://localhost:9000"
	log.Println(result)
	return result, nil
}
