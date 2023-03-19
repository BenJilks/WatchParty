package database

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

type Image struct {
	ID            uint `gorm:"primaryKey;autoIncrement"`
	Title         string
	ThumbnailPath string
	FilePath      string
}

func generateImageThumbnailForFile(thumbnailPath string, imagePath string) string {
	// TODO: Actually generate a smaller thumbnail, instead of just copying
	name := path.Base(imagePath)
	thumbnailFile := path.Join(thumbnailPath, name)

	source, err := os.Open(imagePath)
	if err != nil {
		log.WithError(err).
			Error("Unable to create thumbnail")
		return ""
	}
	defer source.Close()

	destination, err := os.Create(thumbnailFile)
	if err != nil {
		log.WithError(err).
			Error("Unable to create thumbnail")
		return ""
	}
	defer destination.Close()

	if _, err := io.Copy(destination, source); err != nil {
		return ""
	}
	return name
}

func createFileImage(
	db *gorm.DB,
	videoPath string,
	thumbnailPath string,
) (Image, error) {
	thumbnail := generateImageThumbnailForFile(thumbnailPath, videoPath)
	videoFilePath := path.Base(videoPath)
	title := nameFromFile(videoFilePath)
	image := Image{
		Title:         title,
		ThumbnailPath: thumbnail,
		FilePath:      videoFilePath,
	}

	if result := db.Create(&image); result.Error != nil {
		return Image{}, result.Error
	}

	log.WithFields(log.Fields{
		"Title":     title,
		"Thumbnail": thumbnail,
		"File":      videoFilePath,
	}).Info("Registered new file image")

	return image, nil
}

func imageExists(db *gorm.DB, filePath string) (bool, error) {
	var count int64
	result := db.
		Model(&Image{}).
		Where(Image{FilePath: path.Base(filePath)}).
		Count(&count)
	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil
}

func ScanForNewFileImages(db *gorm.DB, imagesPath string, thumbnailsPath string) {
	log.Infof("Scanning '%s' for new file images", imagesPath)
	if err := createDirIfNotExist(thumbnailsPath); err != nil {
		log.WithError(err).
			WithField("path", thumbnailsPath).
			Error("Need path to store thumbnails")
		return
	}

	err := filepath.WalkDir(imagesPath, func(filePath string, info fs.DirEntry, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		exists, err := imageExists(db, filePath)
		if err != nil {
			return err
		}

		if !exists {
			if _, err := createFileImage(db, filePath, thumbnailsPath); err != nil {
				log.WithError(err).
					Warnf("Could not load image file '%s'", filePath)
			}
		}

		return nil
	})

	if err != nil {
		log.WithError(err).
			Error("Error scanning image files")
	}
}
