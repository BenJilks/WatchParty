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
	"strings"
)

const ThumbnailFrameNumber = 20
const ThumbnailScale = 400

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

func videoName(videoName string) string {
	extension := path.Ext(videoName)
	name, _ := strings.CutSuffix(videoName, extension)
	return name
}

func generateThumbnailForFile(thumbnailPath string, inputPath string) (string, error) {
	name := path.Base(inputPath)
	fileName := fmt.Sprintf("%s.jpg", videoName(name))

	outputPath := path.Join(thumbnailPath, fileName)
	if err := readVideoFrame(inputPath, outputPath); err != nil {
		log.Warnf("Unable to generate thumbnail for '%s'\n", inputPath)
		return "", err
	}

	return fileName, nil
}

func createFileVideo(db *gorm.DB, videoPath string, thumbnailPath string) (Video, error) {
	thumbnail, err := generateThumbnailForFile(thumbnailPath, videoPath)
	if err != nil {
		return Video{}, err
	}

	videoFilePath := path.Base(videoPath)
	title := videoName(videoFilePath)
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

func ScanForNewFileVideos(db *gorm.DB, videosPath string, thumbnailsPath string) error {
	log.Infof("Scanning '%s' for new file videos", videosPath)
	if err := createDirIfNotExist(thumbnailsPath); err != nil {
		return err
	}

	return filepath.WalkDir(videosPath, func(filePath string, info fs.DirEntry, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		exists, err := videoExists(db, filePath)
		if err != nil {
			return err
		}

		if !exists {
			if _, err := createFileVideo(db, filePath, thumbnailsPath); err != nil {
				log.WithError(err).
					Warnf("Could not load video file '%s'", filePath)
			}
		}

		return nil
	})
}
