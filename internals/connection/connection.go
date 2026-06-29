package connection

import (
	"hrms/internals/config"
	"hrms/internals/migration"
	"log"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {

	var db *gorm.DB

	connectionKey, err := config.Config()
	if err != nil {
		log.Fatal(err)
		return nil, nil
	}

	maxRetryAttemptsString := config.GetMaxRetryAttempts()
	maxRetryAttempts, err := strconv.Atoi(maxRetryAttemptsString)
	if err != nil {
		log.Fatal(err)
		return nil, nil
	}

	for range maxRetryAttempts {
		db, err = gorm.Open(postgres.Open(connectionKey), &gorm.Config{})
		if err == nil {
			break
		}

		log.Println("Waiting for PostgreSQL...")
		time.Sleep(2 * time.Second)
	}

	log.Println("Connected to PostgreSQL database successfully!")

	err = migration.AutoMigrate(db)
	if err != nil {
		log.Fatal(err)
		return nil, nil
	}

	return db, nil

}

func GetMaxRetryAttempts() any {
	panic("unimplemented")
}
