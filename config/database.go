package config

import (
	"applichic.com/chic_secret/model"
	"database/sql"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB Initialize the database and start automatically the migrations
func InitDB() (*sql.DB, error) {
	//return nil, err
	dbArgs := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		Conf.DatabaseHost, Conf.DatabasePort, Conf.DatabaseUser, Conf.DatabaseName, Conf.DatabasePassword)
	db, err := gorm.Open(postgres.Open(dbArgs), &gorm.Config{})

	// Send the error
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Set the database and migrate the models
	sqlDB.SetMaxIdleConns(Conf.DatabaseMaxConnection)
	DB = db
	db.AutoMigrate(&model.User{}, &model.Token{}, &model.LoginToken{}, &model.Vault{}, &model.Category{},
		&model.Entry{}, &model.CustomField{}, &model.Tag{}, &model.EntryTag{}, &model.VaultUser{})

	return sqlDB, nil
}
