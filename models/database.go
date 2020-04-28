package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"os"
)

var db *gorm.DB // database

func init() {
	fmt.Println("Initializing database...")
	e := godotenv.Load() // Load .env file
	if e != nil {
		fmt.Print(e)
	}

	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")
	dbType := os.Getenv("db_type")

	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s",
		dbHost, username, dbName, password) // Build connection string
	// fmt.Println(dbUri)

	conn, err := gorm.Open(dbType, dbUri)
	if err != nil {
		fmt.Print(err)
		fmt.Printf("Initializing database \t")
		fmt.Println("[FAIL]")
	}

	db = conn
	// db.Debug().AutoMigrate(&Account{}, &Contact{}) // Database migration
	db.Debug().AutoMigrate(&Account{}) // Database migration

	if err == nil {
		fmt.Printf("Initializing database \t")
		fmt.Println("[OK]")
	}
}

// returns a handle to the DB object
func GetDB() *gorm.DB {
	return db
}