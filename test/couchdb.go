package test

import (
	"configurations/databases"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-kivik/kivik/v3"
	"github.com/spf13/viper"
)

// Document represents a simple document for testing
type Document struct {
	ID        string    `json:"_id,omitempty"`
	Rev       string    `json:"_rev,omitempty"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at"`
}

func TestCouchDB() (string, error) {
	// Setup configuration
	viper.SetDefault("couchdb.host", "localhost")
	viper.SetDefault("couchdb.port", 5984)
	viper.SetDefault("couchdb.username", "admin")
	viper.SetDefault("couchdb.password", "password")
	viper.SetDefault("couchdb.protocol", "http")
	viper.SetDefault("couchdb.timeout", 30*time.Second)

	// Initialize CouchDB
	err := databases.InitializeCouchDB()
	if err != nil {
		return "", fmt.Errorf("failed to initialize CouchDB: %w", err)
	}

	log.Println("CouchDB initialized successfully")

	client := databases.GetCouchDBClient()
	ctx := context.Background()

	// List all databases
	log.Println("=== List Databases ===")
	dbNames, err := client.AllDBs(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list databases: %w", err)
	}
	log.Printf("Existing databases: %v", dbNames)

	// Create test database
	dbName := "test_db"
	log.Printf("\n=== Create Database: %s ===", dbName)
	err = client.CreateDB(ctx, dbName)
	if err != nil {
		statusCode := kivik.StatusCode(err)
		if statusCode != 412 {
			return "", fmt.Errorf("failed to create database: %w", err)
		}
		log.Printf("Database %s already exists", dbName)
	} else {
		log.Printf("Database %s created successfully", dbName)
	}

	// Get database instance
	db := client.DB(ctx, dbName)

	// Get database info
	log.Println("\n=== Database Info ===")
	dbInfo, err := db.Stats(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get database info: %w", err)
	}
	log.Printf("Document count: %d", dbInfo.DocCount)

	// Create documents
	log.Println("\n=== Create Documents ===")
	docs := []Document{
		{Name: "Alice", Email: "alice@example.com", Age: 30},
		{Name: "Bob", Email: "bob@example.com", Age: 25},
		{Name: "Charlie", Email: "charlie@example.com", Age: 35},
	}

	var docIDs []string
	var docRevIDs []string
	for i, doc := range docs {
		doc.CreatedAt = time.Now()
		docID, rev, err := db.CreateDoc(ctx, doc)
		if err != nil {
			return "", fmt.Errorf("failed to create document %d: %w", i+1, err)
		}
		docIDs = append(docIDs, docID)
		docRevIDs = append(docRevIDs, rev)
		log.Printf("Created document %d: ID=%s, Rev=%s", i+1, docID, rev)
	}

	// Get document by ID
	log.Println("\n=== Get Document ===")
	var retrievedDoc Document
	if len(docIDs) > 0 {
		row := db.Get(ctx, docIDs[0])
		if err := row.ScanDoc(&retrievedDoc); err != nil {
			return "", fmt.Errorf("failed to get document: %w", err)
		}
		log.Printf("Retrieved: Name=%s, Email=%s, Age=%d", retrievedDoc.Name, retrievedDoc.Email, retrievedDoc.Age)
	}

	// Count documents
	log.Println("\n=== Count Documents ===")
	stats, err := db.Stats(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get database stats: %w", err)
	}
	log.Printf("Total documents: %d", stats.DocCount)

	// Find all documents using _all_docs
	log.Println("\n=== List All Documents ===")
	rows, err := db.AllDocs(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get all docs: %w", err)
	}
	docCount := 0
	for rows.Next() {
		docCount++
	}
	log.Printf("Found %d documents", docCount)

	// Update document
	log.Println("\n=== Update Document ===")
	if len(docIDs) > 0 {
		var docToUpdate Document
		row := db.Get(ctx, docIDs[0])
		if err := row.ScanDoc(&docToUpdate); err != nil {
			return "", fmt.Errorf("failed to get document for update: %w", err)
		}
		docToUpdate.Name = "Alice Updated"
		docToUpdate.Age = 31
		newRev, err := db.Put(ctx, docToUpdate.ID, docToUpdate)
		if err != nil {
			return "", fmt.Errorf("failed to update document: %w", err)
		}
		log.Printf("Updated document: ID=%s, Rev=%s", docToUpdate.ID, newRev)
	}

	// Create a document with specific ID
	log.Println("\n=== Create Document with ID ===")
	specificDoc := Document{
		ID:        "user_john",
		Name:      "John Doe",
		Email:     "john@example.com",
		Age:       28,
		CreatedAt: time.Now(),
	}
	_, err = db.Put(ctx, specificDoc.ID, specificDoc)
	if err != nil {
		return "", fmt.Errorf("failed to create document with ID: %w", err)
	}
	log.Printf("Created document with ID: %s", specificDoc.ID)

	// Get document with specific ID
	var johnDoc Document
	row := db.Get(ctx, "user_john")
	if err := row.ScanDoc(&johnDoc); err != nil {
		return "", fmt.Errorf("failed to get john's document: %w", err)
	}
	log.Printf("Retrieved John: Name=%s, Email=%s", johnDoc.Name, johnDoc.Email)

	// Delete document
	log.Println("\n=== Delete Document ===")
	if len(docIDs) > 0 {
		var docToDelete Document
		row := db.Get(ctx, docIDs[1])
		if err := row.ScanDoc(&docToDelete); err != nil {
			return "", fmt.Errorf("failed to get document for delete: %w", err)
		}
		_, err = db.Delete(ctx, docToDelete.ID, docToDelete.Rev)
		if err != nil {
			return "", fmt.Errorf("failed to delete document: %w", err)
		}
		log.Printf("Deleted document: ID=%s", docToDelete.ID)
	}

	// Final count
	stats, err = db.Stats(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get final stats: %w", err)
	}
	log.Printf("Final document count: %d", stats.DocCount)

	// Clean up - delete database
	log.Printf("\n=== Delete Database: %s ===", dbName)
	err = client.DestroyDB(ctx, dbName)
	if err != nil {
		log.Printf("Warning: failed to delete database: %v", err)
	} else {
		log.Printf("Database %s deleted successfully", dbName)
	}

	// Close connection
	err = databases.CloseCouchDB(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to close CouchDB: %w", err)
	}

	result := fmt.Sprintf("CouchDB test completed!\nTotal operations: 13\nAll operations successful!")
	log.Println(result)

	return result, nil
}

// TestCouchDBWithDSN tests with custom DSN string
func TestCouchDBWithDSN() (string, error) {
	// Setup with DSN string
	dsn := "http://admin:password@localhost:5984"
	viper.SetDefault("couchdb.dsn", dsn)
	viper.SetDefault("couchdb.timeout", 30*time.Second)

	err := databases.InitializeCouchDB()
	if err != nil {
		return "", fmt.Errorf("failed to initialize CouchDB with DSN: %w", err)
	}

	log.Println("CouchDB initialized with DSN successfully")

	client := databases.GetCouchDBClient()
	ctx := context.Background()

	// Simple test - list databases
	dbNames, err := client.AllDBs(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list databases: %w", err)
	}
	log.Printf("Databases: %v", dbNames)

	databases.CloseCouchDB(ctx)

	return "CouchDB DSN test completed successfully!", nil
}
