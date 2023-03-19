package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"path"
)

func Open(filePath string) (*gorm.DB, error) {
	sqliteDB := sqlite.Open(filePath)
	db, err := gorm.Open(sqliteDB, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(Video{}, Image{}); err != nil {
		return nil, err
	}

	return db, nil
}

func nameFromFile(filePath string) string {
	extension := path.Ext(filePath)
	name := filePath[:len(filePath)-len(extension)]
	return name
}
