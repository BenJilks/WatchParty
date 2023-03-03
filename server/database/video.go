package database

import (
	"errors"
	"os"
)

type VideoSourceType = string

const (
	VideoFileSource = VideoSourceType("file")
)

type Video struct {
	ID            uint `gorm:"primaryKey;autoIncrement"`
	Title         string
	ThumbnailPath string

	SourceType    VideoSourceType
	VideoFilePath string
}

func createDirIfNotExist(dirPath string) error {
	_, err := os.Stat(dirPath)
	switch {
	case errors.Is(err, os.ErrNotExist):
		if err := os.MkdirAll(dirPath, os.ModePerm|os.ModeDir); err != nil {
			return err
		}
	case err != nil:
		return err
	}

	return nil
}
