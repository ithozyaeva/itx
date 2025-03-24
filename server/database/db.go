package database

import (
	"fmt"
	"ithozyeva/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDBConnection() error {

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		config.CFG.Database.Host,
		config.CFG.Database.User,
		config.CFG.Database.Password,
		config.CFG.Database.Name,
		config.CFG.Database.Port)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database")
	}

	return nil
}
