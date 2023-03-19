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

func renderImageThumbnail(inputPath string, outputPath string) {
	name := path.Base(outputPath)
	err := ffmpeg.Input(inputPath).
		Filter("scale", ffmpeg.Args{fmt.Sprintf("%d:-1", ThumbnailScale)}).
		Output(outputPath, ffmpeg.KwArgs{"format": "image2", "vcodec": "mjpeg"}).
		OverWriteOutput().
		Run()

	if err != nil {
		log.WithError(err).
			WithFields(log.Fields{"thumbnail": name, "type": "image"}).
			Warnf("Unable to generate thumbnail for '%s'\n", inputPath)
	} else {
		log.WithFields(log.Fields{"thumbnail": name, "type": "image"}).
			Info("Finished creating thumbnail")
	}
}

type ImageThumbnailRequest struct {
	inputPath  string
	outputPath string
}

func imageThumbnailWorker(requests <-chan ImageThumbnailRequest) {
	for request := range requests {
		renderImageThumbnail(request.inputPath, request.outputPath)
	}
}

func startImageThumbnailWorkerPools(count int, requests <-chan ImageThumbnailRequest) {
	for i := 0; i < count; i++ {
		go imageThumbnailWorker(requests)
	}
}

func generateThumbnailForImageFile(
	thumbnailPath string,
	inputPath string,
	thumbnailRequests chan<- ImageThumbnailRequest,
) string {
	name := path.Base(inputPath)
	fileName := fmt.Sprintf("%s.jpg", nameFromFile(name))

	outputPath := path.Join(thumbnailPath, fileName)
	thumbnailRequests <- ImageThumbnailRequest{
		inputPath:  inputPath,
		outputPath: outputPath,
	}

	return fileName
}

func createFileImage(
	db *gorm.DB,
	videoPath string,
	thumbnailPath string,
	thumbnailRequests chan<- ImageThumbnailRequest,
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

func ScanForNewFileImages(db *gorm.DB, imagesPath string, thumbnailsPath string) {
	log.Infof("Scanning '%s' for new file images", imagesPath)
	if err := createDirIfNotExist(thumbnailsPath); err != nil {
		log.WithError(err).
			WithField("path", thumbnailsPath).
			Error("Need path to store thumbnails")
		return
	}

	thumbnailRequests := make(chan ImageThumbnailRequest)
	go startImageThumbnailWorkerPools(ThumbnailWorkerPoolCount, thumbnailRequests)

	err := filepath.WalkDir(imagesPath, func(filePath string, info fs.DirEntry, err error) error {
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
