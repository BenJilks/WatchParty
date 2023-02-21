package main

import (
	"errors"
	"fmt"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"time"
)

const ThumbnailFrameNumber = 20
const ThumbnailScale = 400

type VideoState struct {
	Playing            bool
	Progress           float64
	LastProgressUpdate time.Time
}

type VideoData struct {
	Name          string `json:"name"`
	VideoFile     string `json:"video_file"`
	ThumbnailFile string `json:"thumbnail_file"`
}

func readVideoFrame(path string, outputPath string) error {
	return ffmpeg_go.Input(path).
		Filter("select", ffmpeg_go.Args{fmt.Sprintf("gte(n,%d)", ThumbnailFrameNumber)}).
		Filter("scale", ffmpeg_go.Args{fmt.Sprintf("%d:-1", ThumbnailScale)}).
		Output(outputPath, ffmpeg_go.KwArgs{"frames:v": 1, "format": "image2", "vcodec": "mjpeg"}).
		OverWriteOutput().
		Compile().
		Run()
}

func videoName(videoName string) string {
	extension := path.Ext(videoName)
	nameLength := len(videoName) - len(extension) - 1
	return videoName[:nameLength]
}

func generateThumbnailForFile(thumbnailPath string, inputPath string, info fs.DirEntry) *VideoData {
	if !info.Type().IsRegular() {
		return nil
	}

	name := videoName(info.Name())
	outputPath := path.Join(thumbnailPath, name+".jpg")
	if err := readVideoFrame(inputPath, outputPath); err != nil {
		fmt.Printf("Unable to generate thumbnail for '%s'\n", inputPath)
		return nil
	}

	return &VideoData{
		Name:          name,
		VideoFile:     path.Base(inputPath),
		ThumbnailFile: path.Base(outputPath),
	}
}

func GenerateThumbnails(videoPath string) ([]VideoData, error) {
	thumbnailPath := path.Join(path.Dir(videoPath), "thumbnails")

	_, err := os.Stat(thumbnailPath)
	if errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(thumbnailPath, os.ModePerm|os.ModeDir); err != nil {
			return []VideoData{}, err
		}
	}

	var videos []VideoData
	err = filepath.WalkDir(videoPath, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if videoData := generateThumbnailForFile(thumbnailPath, path, info); videoData != nil {
			videos = append(videos, *videoData)
		}
		return nil
	})

	return videos, err
}
