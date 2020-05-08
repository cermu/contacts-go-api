package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	utl "my-contacts/utils"
	. "os"
)

var db *gorm.DB // database

func init() {
	// fmt.Println("Initializing database...")
	// log.Println("Initializing database...")
	utl.WriteToFile("INFO | Initializing database...")
	e := godotenv.Load() // Load .env file
	if e != nil {
		// fmt.Print(e)
		log.Fatalf("The following error occurred while loading env file: %s", e)
	}

	username := Getenv("db_user")
	password := Getenv("db_pass")
	dbName := Getenv("db_name")
	dbHost := Getenv("db_host")
	dbType := Getenv("db_type")

	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s",
		dbHost, username, dbName, password) // Build connection string
	// fmt.Println(dbUri)

	conn, err := gorm.Open(dbType, dbUri)
	if err != nil {
		// log.Printf("Initializing database \t [FAIL]")
		// log.Printf("The following error occurred while opening DB connection: %s", err)
		utl.WriteToFile(fmt.Sprintf("ERROR | Initializing database \t [FAIL]"))
		utl.WriteToFile(fmt.Sprintf("ERROR | The following error occurred while opening DB connection: %s",
			err))
	}

	db = conn
	db.Debug().AutoMigrate(&Account{}, &Contact{}) // Database migration
	// db.Debug().AutoMigrate(&Account{}) // Database migration

	if err == nil {
		// log.Printf("Initializing database \t [OK]")
		utl.WriteToFile(fmt.Sprintf("INFO | Initializing database \t [OK]"))
	}
}

// returns a handle to the DB object
func GetDB() *gorm.DB {
	return db
}