package databases

import (
	"context"
	"fmt"
	"log"

	"github.com/go-kivik/kivik/v3"
	"github.com/spf13/viper"
)

var CouchDBClient *kivik.Client

func InitializeCouchDB() error {
	host := viper.GetString("couchdb.host")
	port := viper.GetInt("couchdb.port")
	username := viper.GetString("couchdb.username")
	password := viper.GetString("couchdb.password")
	protocol := viper.GetString("couchdb.protocol")

	// Fallback defaults
	if protocol == "" {
		protocol = "http"
	}

	// Construct DSN
	dsn := viper.GetString("couchdb.dsn")
	if dsn == "" {
		if username != "" && password != "" {
			dsn = fmt.Sprintf("%s://%s:%s@%s:%d", protocol, username, password, host, port)
		} else {
			dsn = fmt.Sprintf("%s://%s:%d", protocol, host, port)
		}
	}

	// Create client
	ctx := context.Background()
	client, err := kivik.New("couch", dsn)
	if err != nil {
		return err
	}

	// Test connection - list databases to verify connection
	_, err = client.AllDBs(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to CouchDB: %w", err)
	}

	CouchDBClient = client
	log.Printf("Connected to CouchDB at %s:%d", host, port)
	return nil
}

func GetCouchDBClient() *kivik.Client {
	return CouchDBClient
}

func CloseCouchDB(ctx context.Context) error {
	if CouchDBClient != nil {
		if err := CouchDBClient.Close(ctx); err != nil {
			return err
		}
		log.Println("CouchDB connection closed")
	}
	return nil
}

func NewCouchDBClient() (*kivik.Client, error) {
	host := viper.GetString("couchdb.host")
	port := viper.GetInt("couchdb.port")
	username := viper.GetString("couchdb.username")
	password := viper.GetString("couchdb.password")
	protocol := viper.GetString("couchdb.protocol")

	if protocol == "" {
		protocol = "http"
	}

	dsn := viper.GetString("couchdb.dsn")
	if dsn == "" {
		if username != "" && password != "" {
			dsn = fmt.Sprintf("%s://%s:%s@%s:%d", protocol, username, password, host, port)
		} else {
			dsn = fmt.Sprintf("%s://%s:%d", protocol, host, port)
		}
	}

	client, err := kivik.New("couch", dsn)
	if err != nil {
		return nil, err
	}

	return client, nil
}
