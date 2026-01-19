package test

import (
	"configurations/databases"
	"fmt"
	"log"
	"time"

	"github.com/gocql/gocql"
	"github.com/spf13/viper"
)

func TestCassandra() (string, error) {
	// Setup configuration
	viper.SetDefault("cassandra.hosts", []string{"localhost"})
	viper.SetDefault("cassandra.port", 9042)
	viper.SetDefault("cassandra.keyspace", "system")
	viper.SetDefault("cassandra.num_conns", 2)
	viper.SetDefault("cassandra.timeout", 30*time.Second)
	viper.SetDefault("cassandra.consistency", "quorum")

	// Initialize Cassandra
	err := databases.InitializeCassandra()
	if err != nil {
		return "", fmt.Errorf("failed to initialize Cassandra: %w", err)
	}

	log.Println("Cassandra initialized successfully")

	session := databases.GetCassandraSession()

	// Test query: Get cluster info
	var releaseVersion string
	err = session.Query("SELECT release_version FROM system.local").Consistency(gocql.One).Scan(&releaseVersion)
	if err != nil {
		return "", fmt.Errorf("failed to query system.local: %w", err)
	}
	log.Printf("Cassandra release version: %s", releaseVersion)

	// Test query: Get cluster name
	var clusterName string
	err = session.Query("SELECT cluster_name FROM system.local").Consistency(gocql.One).Scan(&clusterName)
	if err != nil {
		return "", fmt.Errorf("failed to query cluster name: %w", err)
	}
	log.Printf("Cluster name: %s", clusterName)

	// Test query: Get all peers
	hosts, err := session.Query("SELECT peer FROM system.peers").Iter().SliceMap()
	if err != nil {
		log.Printf("Warning: failed to query peers (might be single node): %v", err)
	} else {
		log.Printf("Found %d peer(s)", len(hosts))
	}

	// Test query: Get data center
	var dataCenter string
	err = session.Query("SELECT data_center FROM system.local").Consistency(gocql.One).Scan(&dataCenter)
	if err != nil {
		log.Printf("Warning: failed to query data center: %v", err)
	} else {
		log.Printf("Data center: %s", dataCenter)
	}

	// Test: Create keyspace and table
	keyspaceName := "test_keyspace"
	err = session.Query(fmt.Sprintf(
		"CREATE KEYSPACE IF NOT EXISTS %s WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}",
		keyspaceName,
	)).Exec()
	if err != nil {
		return "", fmt.Errorf("failed to create keyspace: %w", err)
	}
	log.Printf("Created keyspace: %s", keyspaceName)

	// Close session and reconnect with new keyspace
	err = databases.CloseCassandra()
	if err != nil {
		return "", err
	}
	viper.Set("cassandra.keyspace", keyspaceName)
	err = databases.InitializeCassandra()
	if err != nil {
		return "", fmt.Errorf("failed to reconnect with new keyspace: %w", err)
	}
	session = databases.GetCassandraSession()

	// Create table
	err = session.Query(
		"CREATE TABLE IF NOT EXISTS users (id UUID PRIMARY KEY, name text, email text, created_at timestamp)",
	).Exec()
	if err != nil {
		return "", fmt.Errorf("failed to create table: %w", err)
	}
	log.Println("Created table: users")

	// Insert test data
	id := gocql.TimeUUID()
	err = session.Query(
		"INSERT INTO users (id, name, email, created_at) VALUES (?, ?, ?, toTimestamp(now()))",
		id, "Test User", "test@example.com",
	).Exec()
	if err != nil {
		return "", fmt.Errorf("failed to insert data: %w", err)
	}
	log.Printf("Inserted test data with ID: %s", id)

	// Query back the data
	var name, email string
	err = session.Query(
		"SELECT name, email FROM users WHERE id = ?",
		id,
	).Consistency(gocql.One).Scan(&name, &email)
	if err != nil {
		return "", fmt.Errorf("failed to query data: %w", err)
	}
	log.Printf("Retrieved: name=%s, email=%s", name, email)

	// Count rows
	var count int
	err = session.Query("SELECT COUNT(*) FROM users").Consistency(gocql.One).Scan(&count)
	if err != nil {
		return "", fmt.Errorf("failed to count rows: %w", err)
	}
	log.Printf("Total rows in users table: %d", count)

	// Clean up
	err = session.Query(fmt.Sprintf("DROP KEYSPACE %s", keyspaceName)).Exec()
	if err != nil {
		log.Printf("Warning: failed to drop keyspace: %v", err)
	} else {
		log.Printf("Dropped keyspace: %s", keyspaceName)
	}

	err = databases.CloseCassandra()
	if err != nil {
		return "", err
	}

	result := fmt.Sprintf("Cassandra test completed!\nVersion: %s\nCluster: %s\nAll operations successful!",
		releaseVersion, clusterName)
	log.Println(result)

	return result, nil
}
