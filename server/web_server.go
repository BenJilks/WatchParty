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
	"sync"
	"time"
)

const TempDirectoryName = "watch-party"

type CachedFile struct {
	Time         time.Time
	BeingWritten bool
	Size         int64
}

type GzipFileCache struct {
	tempDirectory string
	cache         map[string]CachedFile
	mutex         sync.Mutex
}

type PathType int

const (
	PathTypeNothing = PathType(iota)
	PathTypeFile
	PathTypeDirectory
)

type PathDescription struct {
	pathType PathType
	rawPath  string

	contentType  *string
	size         *int64
	lastModified *time.Time
}

type DoubleWriter struct {
	first  io.Writer
	second io.Writer
}

func (doubleWriter DoubleWriter) Write(data []byte) (int, error) {
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

func readPathDescription(filePath string) PathDescription {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return PathDescription{
			pathType: PathTypeNothing,
			rawPath:  filePath,
		}
	}

	lastModified := fileInfo.ModTime()
	if fileInfo.IsDir() {
		return PathDescription{
			pathType:     PathTypeDirectory,
			rawPath:      filePath,
			lastModified: &lastModified,
		}
	}

	contentType := mime.TypeByExtension(path.Ext(filePath))
	size := fileInfo.Size()
	return PathDescription{
		pathType:     PathTypeFile,
		rawPath:      filePath,
		contentType:  &contentType,
		size:         &size,
		lastModified: &lastModified,
	}
}

func gzipAndServeFile(filePath string, gzippedFilePath string, response http.ResponseWriter) (int64, error) {
	originalFile, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer originalFile.Close()

	gzippedFile, err := os.Create(gzippedFilePath)
	if err != nil {
		return 0, err
	}
	defer gzippedFile.Close()

	writer := gzip.NewWriter(&DoubleWriter{
		first:  gzippedFile,
		second: response,
	})
	defer writer.Close()

	response.Header().Set("Content-Encoding", "gzip")
	return io.Copy(writer, originalFile)
}

func serveCachedGzippedFile(response http.ResponseWriter, filePath string, size int64) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	response.Header().Set("Content-Length", fmt.Sprint(size))
	response.Header().Set("Content-Encoding", "gzip")
	_, err = io.Copy(response, file)
	return err
}

func (fileCache *GzipFileCache) getCachedGzippedFile(description PathDescription) (string, *CachedFile) {
	cacheName := strings.ReplaceAll(description.rawPath, "/", "_")
	cacheName = strings.ReplaceAll(cacheName, ".", "_")
	gzippedFilePath := path.Join(fileCache.tempDirectory, cacheName+".gz")

	fileCache.mutex.Lock()
	cachedFile, inCache := fileCache.cache[description.rawPath]
	fileCache.mutex.Unlock()

	if inCache && !description.lastModified.After(cachedFile.Time) {
		return gzippedFilePath, &cachedFile
	}

	return gzippedFilePath, nil
}

func (fileCache *GzipFileCache) cacheAndServeFile(
	response http.ResponseWriter,
	description PathDescription,
	gzippedFilePath string,
) error {
	log.WithField("file", description.rawPath).
		Info("Updating gzip cache")

	fileCache.mutex.Lock()
	fileCache.cache[description.rawPath] = CachedFile{
		Time:         *description.lastModified,
		BeingWritten: true,
	}
	fileCache.mutex.Unlock()

	size, err := gzipAndServeFile(description.rawPath, gzippedFilePath, response)
	if err != nil {
		delete(fileCache.cache, description.rawPath)
		return err
	}

	fileCache.mutex.Lock()
	fileCache.cache[description.rawPath] = CachedFile{
		Time:         *description.lastModified,
		BeingWritten: false,
		Size:         size,
	}
	fileCache.mutex.Unlock()
	return nil
}

func (fileCache *GzipFileCache) serveGzipFile(
	response http.ResponseWriter,
	description PathDescription,
) error {
	if description.pathType != PathTypeFile {
		return errors.New("file doesn't exist")
	}

	response.Header().Set("Content-Type", *description.contentType)
	gzippedFilePath, cachedFile := fileCache.getCachedGzippedFile(description)
	if cachedFile != nil {
		if cachedFile.BeingWritten {
			return errors.New("file currently being cached")
		}
		return serveCachedGzippedFile(response, gzippedFilePath, cachedFile.Size)
	}

	return fileCache.cacheAndServeFile(
		response, description, gzippedFilePath)
}

func (fileCache *GzipFileCache) serveFile(
	response http.ResponseWriter,
	request *http.Request,
	description PathDescription,
) {
	acceptEncoding := request.Header.Get("Accept-Encoding")
	if !strings.Contains(acceptEncoding, "gzip") {
		http.ServeFile(response, request, description.rawPath)
		return
	}

	// NOTE: Don't gzip video content.
	mimeType := mime.TypeByExtension(path.Ext(description.rawPath))
	if len(mimeType) >= 5 && mimeType[:5] == "video" {
		http.ServeFile(response, request, description.rawPath)
		return
	}

	if err := fileCache.serveGzipFile(response, description); err != nil {
		log.WithError(err).
			WithField("file", description.rawPath).
			Error("Could not serve file gzipped")
		http.ServeFile(response, request, description.rawPath)
	}
}

func firstValidIndexPath(directoryPath string) PathDescription {
	validIndexFiles := []string{
		"index.html",
		"index.htm",
	}

	for _, indexFile := range validIndexFiles {
		indexPath := path.Join(directoryPath, indexFile)
		description := readPathDescription(indexPath)
		if description.pathType == PathTypeFile {
			return description
		}
	}

	return PathDescription{
		pathType: PathTypeNothing,
	}
}

func (fileCache *GzipFileCache) serveDirectory(
	response http.ResponseWriter,
	request *http.Request,
	directoryPath string,
) {
	if description := firstValidIndexPath(directoryPath); description.pathType != PathTypeNothing {
		fileCache.serveFile(response, request, description)
		return
	}

	http.ServeFile(response, request, directoryPath)
}

func WebHandler(staticPath string) http.HandlerFunc {
	tempDirectory := path.Join(os.TempDir(), TempDirectoryName)
	_ = os.MkdirAll(tempDirectory, os.ModeDir|os.ModePerm)

	log.WithField("cache-path", tempDirectory).
		Info("Using gzip cache")

	fileCache := GzipFileCache{
		tempDirectory: tempDirectory,
		cache:         map[string]CachedFile{},
		mutex:         sync.Mutex{},
	}

	return func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		response.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")

		url := request.URL.Path
		filePath := path.Join(staticPath, url)
		description := readPathDescription(filePath)
		switch description.pathType {
		case PathTypeNothing:
			http.ServeFile(response, request, filePath)
		case PathTypeFile:
			fileCache.serveFile(response, request, description)
		case PathTypeDirectory:
			fileCache.serveDirectory(response, request, filePath)
		}
	}
}
