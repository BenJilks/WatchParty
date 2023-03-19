package database

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"gorm.io/gorm"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

const ThumbnailFrameNumber = 24 * 10

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

func readVideoFrame(path string, outputPath string) error {
	return ffmpeg.Input(path).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", ThumbnailFrameNumber)}).
		Filter("scale", ffmpeg.Args{fmt.Sprintf("%d:-1", ThumbnailScale)}).
		Output(outputPath, ffmpeg.KwArgs{"frames:v": 1, "format": "image2", "vcodec": "mjpeg"}).
		OverWriteOutput().
		Run()
}

func renderVideoThumbnail(inputPath string, outputPath string) {
	name := path.Base(outputPath)
	if err := readVideoFrame(inputPath, outputPath); err != nil {
		log.WithError(err).
			WithField("thumbnail", name).
			Warnf("Unable to generate thumbnail for '%s'\n", inputPath)
	} else {
		log.WithField("thumbnail", name).
			Info("Finished creating thumbnail")
	}
}

type VideoThumbnailRequest struct {
	inputPath  string
	outputPath string
}

func videoThumbnailWorker(requests <-chan VideoThumbnailRequest) {
	for request := range requests {
		renderVideoThumbnail(request.inputPath, request.outputPath)
	}
}

func startVideoThumbnailWorkerPools(count int, requests <-chan VideoThumbnailRequest) {
	for i := 0; i < count; i++ {
		go videoThumbnailWorker(requests)
	}
}

func generateThumbnailForVideoFile(
	thumbnailPath string,
	inputPath string,
	thumbnailRequests chan<- VideoThumbnailRequest,
) string {
	name := path.Base(inputPath)
	fileName := fmt.Sprintf("%s.jpg", nameFromFile(name))

	outputPath := path.Join(thumbnailPath, fileName)
	thumbnailRequests <- VideoThumbnailRequest{
		inputPath:  inputPath,
		outputPath: outputPath,
	}

	return fileName
}

func createFileVideo(
	db *gorm.DB,
	videoPath string,
	thumbnailPath string,
	thumbnailRequests chan<- VideoThumbnailRequest,
) (Video, error) {
	thumbnail := generateThumbnailForVideoFile(thumbnailPath, videoPath, thumbnailRequests)
	videoFilePath := path.Base(videoPath)
	title := nameFromFile(videoFilePath)
	video := Video{
		Title:         title,
		ThumbnailPath: thumbnail,
		SourceType:    VideoFileSource,
		VideoFilePath: videoFilePath,
	}

	if result := db.Create(&video); result.Error != nil {
		return Video{}, result.Error
	}

	log.WithFields(log.Fields{
		"Title":     title,
		"Thumbnail": thumbnail,
		"File":      videoFilePath,
	}).Info("Registered new file video")

	return video, nil
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

func videoExists(db *gorm.DB, filePath string) (bool, error) {
	var count int64
	result := db.
		Model(&Video{}).
		Where(Video{VideoFilePath: path.Base(filePath)}).
		Count(&count)
	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil
}

func ScanForNewFileVideos(db *gorm.DB, videosPath string, thumbnailsPath string) {
	log.Infof("Scanning '%s' for new file videos", videosPath)
	if err := createDirIfNotExist(thumbnailsPath); err != nil {
		log.WithError(err).
			WithField("path", thumbnailsPath).
			Error("Need path to store thumbnails")
		return
	}

	thumbnailRequests := make(chan VideoThumbnailRequest)
	go startVideoThumbnailWorkerPools(ThumbnailWorkerPoolCount, thumbnailRequests)

	err := filepath.WalkDir(videosPath, func(filePath string, info fs.DirEntry, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		exists, err := videoExists(db, filePath)
		if err != nil {
			return err
		}

		if !exists {
			if _, err := createFileVideo(db, filePath, thumbnailsPath, thumbnailRequests); err != nil {
				log.WithError(err).
					Warnf("Could not load video file '%s'", filePath)
			}
		}

		return nil
	})

	if err != nil {
		log.WithError(err).
			Error("Error scanning video files")
	}
}
