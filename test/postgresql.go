package test

import (
	"configurations/databases"
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

type User struct {
	ID        uint      `gorm:"primarykey"`
	Name      string    `gorm:"size:255;not null"`
	Email     string    `gorm:"size:255;uniqueIndex;not null"`
	CreatedAt time.Time `gorm:"autoUpdateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func TestPostgreSQL() (string, error) {
	// Setup configuration
	viper.SetDefault("postgresql.host", "localhost")
	viper.SetDefault("postgresql.port", 5432)
	viper.SetDefault("postgresql.user", "postgres")
	viper.SetDefault("postgresql.password", "postgres")
	viper.SetDefault("postgresql.dbname", "exam_db")
	viper.SetDefault("postgresql.sslmode", "disable")
	viper.SetDefault("postgresql.log_level", "info")
	viper.SetDefault("postgresql.max_idle_conns", 10)
	viper.SetDefault("postgresql.max_open_conns", 100)
	viper.SetDefault("postgresql.conn_max_lifetime", "1h")

	// Initialize PostgreSQL
	err := databases.InitializePostgersql()
	if err != nil {
		return "", fmt.Errorf("failed to initialize PostgreSQL: %w", err)
	}

	log.Println("PostgreSQL initialized successfully")

	db := databases.GetPostgresDB()

	// Get database version
	var version string
	err = db.Raw("SELECT version()").Scan(&version).Error
	if err != nil {
		return "", fmt.Errorf("failed to query version: %w", err)
	}
	log.Printf("PostgreSQL version: %s", version)

	// Get current database
	var currentDB string
	err = db.Raw("SELECT current_database()").Scan(&currentDB).Error
	if err != nil {
		return "", fmt.Errorf("failed to query current database: %w", err)
	}
	log.Printf("Current database: %s", currentDB)

	// Get connection stats
	sqlDB, _ := db.DB()
	stats := sqlDB.Stats()
	log.Printf("Connection stats - Open: %d, InUse: %d, Idle: %d",
		stats.OpenConnections, stats.InUse, stats.Idle)

	// Auto migrate User table
	err = db.AutoMigrate(&User{})
	if err != nil {
		return "", fmt.Errorf("failed to migrate User table: %w", err)
	}
	log.Println("User table migrated successfully")

	// Create test users
	users := []User{
		{Name: "Alice", Email: "alice@example.com"},
		{Name: "Bob", Email: "bob@example.com"},
		{Name: "Charlie", Email: "charlie@example.com"},
	}

	for i, user := range users {
		err = db.Create(&user).Error
		if err != nil {
			return "", fmt.Errorf("failed to create user %d: %w", i+1, err)
		}
		log.Printf("Created user: %s (%s)", user.Name, user.Email)
	}

	// Count users
	var count int64
	err = db.Model(&User{}).Count(&count).Error
	if err != nil {
		return "", fmt.Errorf("failed to count users: %w", err)
	}
	log.Printf("Total users: %d", count)

	// Find all users
	var retrievedUsers []User
	err = db.Find(&retrievedUsers).Error
	if err != nil {
		return "", fmt.Errorf("failed to find users: %w", err)
	}
	log.Printf("Retrieved %d users", len(retrievedUsers))

	// Find user by email
	var alice User
	err = db.Where("email = ?", "alice@example.com").First(&alice).Error
	if err != nil {
		return "", fmt.Errorf("failed to find alice: %w", err)
	}
	log.Printf("Found by email: %s (%s)", alice.Name, alice.Email)

	// Update user
	alice.Name = "Alice Updated"
	err = db.Save(&alice).Error
	if err != nil {
		return "", fmt.Errorf("failed to update alice: %w", err)
	}
	log.Printf("Updated user: %s", alice.Name)

	// Delete user
	err = db.Where("email = ?", "bob@example.com").Delete(&User{}).Error
	if err != nil {
		return "", fmt.Errorf("failed to delete bob: %w", err)
	}
	log.Println("Deleted user: bob@example.com")

	// Final count
	err = db.Model(&User{}).Count(&count).Error
	if err != nil {
		return "", fmt.Errorf("failed to count users after delete: %w", err)
	}
	log.Printf("Total users after delete: %d", count)

	// Clean up - drop table
	err = db.Migrator().DropTable(&User{})
	if err != nil {
		log.Printf("Warning: failed to drop table: %v", err)
	} else {
		log.Println("Dropped User table")
	}

	// Close connection
	err = databases.ClosePostgres()
	if err != nil {
		return "", fmt.Errorf("failed to close PostgreSQL: %w", err)
	}

	result := fmt.Sprintf("PostgreSQL test completed!\nDatabase: %s\nTotal operations: 10\nAll operations successful!",
		currentDB)
	log.Println(result)

	return result, nil
}

// TestPostgreSQLWithDSN tests with custom DSN string
func TestPostgreSQLWithDSN() (string, error) {
	// Setup with DSN string
	dsn := "host=localhost port=5432 user=postgres password=postgres dbname=exam_db sslmode=disable"
	viper.SetDefault("postgresql.dsn", dsn)
	viper.SetDefault("postgresql.log_level", "info")

	err := databases.InitializePostgersql()
	if err != nil {
		return "", fmt.Errorf("failed to initialize PostgreSQL with DSN: %w", err)
	}

	log.Println("PostgreSQL initialized with DSN successfully")

	db := databases.GetPostgresDB()

	// Simple test query
	var result string
	err = db.Raw("SELECT 'Connection successful!' as message").Scan(&result).Error
	if err != nil {
		return "", fmt.Errorf("failed to execute test query: %w", err)
	}
	log.Println(result)

	databases.ClosePostgres()

	return "PostgreSQL DSN test completed successfully!", nil
}
