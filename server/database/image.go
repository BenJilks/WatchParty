package database

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"gorm.io/gorm"
	"io/fs"
	"path"
	"path/filepath"
)

type Image struct {
	ID            uint `gorm:"primaryKey;autoIncrement"`
	Title         string
	ThumbnailPath string
	FilePath      string
}

func renderImageThumbnail(inputPath string, outputPath string) *ffmpeg.Stream {
	return ffmpeg.Input(inputPath).
		Filter("scale", ffmpeg.Args{fmt.Sprintf("%d:-1", ThumbnailScale)}).
		Output(outputPath, ffmpeg.KwArgs{"format": "image2", "vcodec": "mjpeg"}).
		OverWriteOutput()
}

func generateThumbnailForImageFile(
	thumbnailPath string,
	inputPath string,
	thumbnailRequests chan<- FFMPegRequest,
) string {
	name := path.Base(inputPath)
	fileName := fmt.Sprintf("%s.jpg", nameFromFile(name))

	outputPath := path.Join(thumbnailPath, fileName)
	thumbnailRequests <- FFMPegRequest{
		inputPath:  inputPath,
		outputPath: outputPath,
		stream:     renderImageThumbnail(inputPath, outputPath),
		typeName:   "image",
	}

	return fileName
}

func createFileImage(
	db *gorm.DB,
	videoPath string,
	thumbnailPath string,
	thumbnailRequests chan<- FFMPegRequest,
) (Image, error) {
	thumbnail := generateThumbnailForImageFile(thumbnailPath, videoPath, thumbnailRequests)
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

func ScanForNewFileImages(
	db *gorm.DB,
	imagesPath string,
	thumbnailsPath string,
	thumbnailRequests chan<- FFMPegRequest,
) {
	log.Infof("Scanning '%s' for new file images", imagesPath)
	if err := createDirIfNotExist(thumbnailsPath); err != nil {
		log.WithError(err).
			WithField("path", thumbnailsPath).
			Error("Need path to store thumbnails")
		return
	}

	err := filepath.WalkDir(imagesPath+"/", func(filePath string, info fs.DirEntry, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		exists, err := imageExists(db, filePath)
		if err != nil {
			return err
		}

		if !exists {
			if _, err := createFileImage(db, filePath, thumbnailsPath, thumbnailRequests); err != nil {
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
