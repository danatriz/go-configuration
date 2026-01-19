package databases

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var PostgreDB *gorm.DB

func InitializePostgersql() error {
	host := viper.GetString("postgresql.host")
	port := viper.GetInt("postgresql.port")
	user := viper.GetString("postgresql.user")
	password := viper.GetString("postgresql.password")
	dbname := viper.GetString("postgresql.dbname")
	sslmode := viper.GetString("postgresql.sslmode")

	dsn := getPostgresDSN(host, port, user, password, dbname, sslmode)

	logLevel := getPostgresLogLevel(viper.GetString("postgresql.log_level"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(logLevel),
	})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// Set connection pool settings
	maxIdleConns := viper.GetInt("postgresql.max_idle_conns")
	maxOpenConns := viper.GetInt("postgresql.max_open_conns")
	connMaxLifetime := viper.GetDuration("postgresql.conn_max_lifetime")

	if maxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(maxIdleConns)
	}
	if maxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(maxOpenConns)
	}
	if connMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(connMaxLifetime)
	}

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return err
	}

	PostgreDB = db
	log.Printf("Connected to PostgreSQL at %s:%d, database: %s", host, port, dbname)
	return nil
}

func GetPostgresDB() *gorm.DB {
	return PostgreDB
}

func ClosePostgres() error {
	if PostgreDB != nil {
		sqlDB, err := PostgreDB.DB()
		if err != nil {
			return err
		}
		if err := sqlDB.Close(); err != nil {
			return err
		}
		log.Println("PostgreSQL connection closed")
	}
	return nil
}

func getPostgresDSN(host string, port int, user, password, dbname, sslmode string) string {
	// If dsn is provided, use it
	if dsn := viper.GetString("postgresql.dsn"); dsn != "" {
		return dsn
	}
	// Otherwise construct from components
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)
}

func getPostgresLogLevel(level string) gormlogger.LogLevel {
	switch level {
	case "silent":
		return gormlogger.Silent
	case "error":
		return gormlogger.Error
	case "warn":
		return gormlogger.Warn
	case "info":
		return gormlogger.Info
	default:
		return gormlogger.Info
	}
}
