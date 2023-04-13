package main

import (
	"compress/gzip"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"mime"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

const TempDirectoryName = "watch-party"

type GzipFileCache struct {
	tempDirectory string
	cache         map[string]time.Time
}

type FileDescription struct {
	contentType  string
	size         *int64
	lastModified *time.Time
}

type PathType int

const (
	PathTypeNothing = PathType(iota)
	PathTypeFile
	PathTypeDirectory
)

func readFileDescription(filePath string) FileDescription {
	contentType := mime.TypeByExtension(path.Ext(filePath))

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return FileDescription{
			contentType: contentType,
		}
	}

	size := fileInfo.Size()
	lastModified := fileInfo.ModTime()
	return FileDescription{
		contentType:  contentType,
		size:         &size,
		lastModified: &lastModified,
	}
}

type DoubleWriter struct {
	first  io.Writer
	second io.Writer
}

func (doubleWriter *DoubleWriter) Write(data []byte) (int, error) {
	firstCount, err := doubleWriter.first.Write(data)
	if err != nil {
		return firstCount, err
	}

	secondCount, err := doubleWriter.second.Write(data[:firstCount])
	if err != nil {
		return secondCount, err
	}

	return firstCount, nil
}

func gzipAndServeFile(filePath string, gzippedFilePath string, response http.ResponseWriter) error {
	originalFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer originalFile.Close()

	gzippedFile, err := os.Create(gzippedFilePath)
	if err != nil {
		return err
	}
	defer gzippedFile.Close()

	writer := gzip.NewWriter(&DoubleWriter{
		first:  gzippedFile,
		second: response,
	})
	defer writer.Close()

	response.Header().Set("Content-Encoding", "gzip")
	_, err = io.Copy(writer, originalFile)
	return err
}

func serveCachedGzippedFile(response http.ResponseWriter, filePath string) error {
	if fileInfo, err := os.Stat(filePath); err == nil {
		response.Header().Set("Content-Length", fmt.Sprint(fileInfo.Size()))
	}

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	response.Header().Set("Content-Encoding", "gzip")
	_, err = io.Copy(response, file)
	return err
}

func (fileCache *GzipFileCache) getCachedGzippedFile(description FileDescription, filePath string) (string, bool) {
	cacheName := strings.ReplaceAll(filePath, "/", "_")
	cacheName = strings.ReplaceAll(cacheName, ".", "_")
	gzippedFilePath := path.Join(fileCache.tempDirectory, cacheName+".gz")

	cacheLastModified, inCache := fileCache.cache[filePath]
	if inCache && !description.lastModified.After(cacheLastModified) {
		return gzippedFilePath, true
	}

	return gzippedFilePath, false
}

func (fileCache *GzipFileCache) cacheAndServeFile(
	response http.ResponseWriter,
	description FileDescription,
	filePath string,
	gzippedFilePath string,
) error {
	log.WithField("file", filePath).
		Info("Updating gzip cache")

	if err := gzipAndServeFile(filePath, gzippedFilePath, response); err != nil {
		return err
	}

	fileCache.cache[filePath] = *description.lastModified
	return nil
}

func (fileCache *GzipFileCache) serveGzipFile(response http.ResponseWriter, filePath string) error {
	description := readFileDescription(filePath)
	if description.lastModified == nil {
		return errors.New("could not stat file")
	}

	response.Header().Set("Content-Type", description.contentType)
	gzippedFilePath, cached := fileCache.getCachedGzippedFile(description, filePath)
	if cached {
		return serveCachedGzippedFile(response, gzippedFilePath)
	}

	return fileCache.cacheAndServeFile(
		response, description, filePath, gzippedFilePath)
}

func (fileCache *GzipFileCache) serveFile(response http.ResponseWriter, request *http.Request, filePath string) {
	acceptEncoding := request.Header.Get("Accept-Encoding")
	if !strings.Contains(acceptEncoding, "gzip") {
		http.ServeFile(response, request, filePath)
		return
	}

	if err := fileCache.serveGzipFile(response, filePath); err != nil {
		log.WithError(err).
			WithField("file", filePath).
			Error("Could not serve file gzipped")
		http.ServeFile(response, request, filePath)
	}

}

func getPathType(filePath string) PathType {
	info, err := os.Stat(filePath)
	if err != nil {
		return PathTypeNothing
	}

	if info.IsDir() {
		return PathTypeDirectory
	} else {
		return PathTypeFile
	}
}

func firstValidIndexPath(directoryPath string) *string {
	validIndexFiles := []string{
		"index.html",
		"index.htm",
	}

	for _, indexFile := range validIndexFiles {
		indexPath := path.Join(directoryPath, indexFile)
		if getPathType(indexPath) == PathTypeFile {
			return &indexPath
		}
	}

	return nil
}

func (fileCache *GzipFileCache) serveDirectory(
	response http.ResponseWriter,
	request *http.Request,
	directoryPath string,
) {
	if indexPath := firstValidIndexPath(directoryPath); indexPath != nil {
		fileCache.serveFile(response, request, *indexPath)
		return
	}

	http.ServeFile(response, request, directoryPath)
}

func WebHandler(staticPath string) http.HandlerFunc {
	tempDirectory := path.Join(os.TempDir(), TempDirectoryName)
	_ = os.Mkdir(tempDirectory, os.ModeDir|os.ModePerm)

	log.WithField("cache-path", tempDirectory).
		Info("Using gzip cache")

	fileCache := GzipFileCache{
		tempDirectory: tempDirectory,
		cache:         map[string]time.Time{},
	}

	return func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		response.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")

		url := request.URL.Path
		filePath := path.Join(staticPath, url)
		switch getPathType(filePath) {
		case PathTypeNothing:
			http.ServeFile(response, request, filePath)
		case PathTypeFile:
			fileCache.serveFile(response, request, filePath)
		case PathTypeDirectory:
			fileCache.serveDirectory(response, request, filePath)
		}
	}
}
