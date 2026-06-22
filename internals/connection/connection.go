package connection

import (
	"fmt"
	"hrms/internals/config"
	"hrms/internals/migration"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {

	connectionKey, err := config.Config()
	if err != nil {
		log.Fatal(err)
		return nil, nil
	}

	db, err := gorm.Open(postgres.Open(connectionKey), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("Failed connecting to DB : %v", err)
	}

	err = migration.AutoMigrate(db)
	if err != nil {
		log.Fatal(err)
		return nil, nil
	}

	return db, nil

}
