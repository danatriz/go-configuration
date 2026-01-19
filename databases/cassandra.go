package databases

import (
	"log"

	"github.com/gocql/gocql"
	"github.com/spf13/viper"
)

var CassandraSession *gocql.Session

func InitializeCassandra() error {
	// Get configuration from viper
	hosts := viper.GetStringSlice("cassandra.hosts")
	keyspace := viper.GetString("cassandra.keyspace")
	port := viper.GetInt("cassandra.port")
	conns := viper.GetInt("cassandra.num_conns")
	timeout := viper.GetDuration("cassandra.timeout")

	// Create cluster configuration
	cluster := gocql.NewCluster(hosts...)
	cluster.Keyspace = keyspace
	cluster.Port = port
	cluster.NumConns = conns
	cluster.Timeout = timeout

	// Consistency level
	consistency := viper.GetString("cassandra.consistency")
	switch consistency {
	case "one":
		cluster.Consistency = gocql.One
	case "quorum":
		cluster.Consistency = gocql.Quorum
	case "all":
		cluster.Consistency = gocql.All
	case "localQuorum":
		cluster.Consistency = gocql.LocalQuorum
	default:
		cluster.Consistency = gocql.Quorum
	}

	// Create session
	session, err := cluster.CreateSession()
	if err != nil {
		return err
	}

	CassandraSession = session
	log.Printf("Connected to Cassandra cluster at %v, keyspace: %s", hosts, keyspace)
	return nil
}

func GetCassandraSession() *gocql.Session {
	return CassandraSession
}

func CloseCassandra() error {
	if CassandraSession != nil {
		CassandraSession.Close()
		log.Println("Cassandra session closed")
	}
	return nil
}

func NewCassandraSession() (*gocql.Session, error) {
	hosts := viper.GetStringSlice("cassandra.hosts")
	keyspace := viper.GetString("cassandra.keyspace")
	port := viper.GetInt("cassandra.port")
	conns := viper.GetInt("cassandra.num_conns")
	timeout := viper.GetDuration("cassandra.timeout")

	cluster := gocql.NewCluster(hosts...)
	cluster.Keyspace = keyspace
	cluster.Port = port
	cluster.NumConns = conns
	cluster.Timeout = timeout

	consistency := viper.GetString("cassandra.consistency")
	switch consistency {
	case "one":
		cluster.Consistency = gocql.One
	case "quorum":
		cluster.Consistency = gocql.Quorum
	case "all":
		cluster.Consistency = gocql.All
	case "localQuorum":
		cluster.Consistency = gocql.LocalQuorum
	default:
		cluster.Consistency = gocql.Quorum
	}

	return cluster.CreateSession()
}
