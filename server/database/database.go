package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Open(filePath string) (*gorm.DB, error) {
	sqliteDB := sqlite.Open(filePath)
	db, err := gorm.Open(sqliteDB, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(Video{}); err != nil {
		return nil, err
	}

	return db, nil
}
