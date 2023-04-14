package main

import (
	"flag"
	webserver "github.com/benjilks/tinywebserver"
	log "github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"gorm.io/gorm"
	"runtime"
	"watch-party/database"
)

type Config struct {
	LogLevel string

	DatabasePath   string
	VideosPath     string
	ImagesPath     string
	ThumbnailsPath string

	webserver.WebServerConfig
}

func defaultConfig() Config {
	webserverConfig := webserver.DefaultWebServerConfig()
	webserverConfig.StaticFilesPath = DefaultStaticFilesPath

	return Config{
		LogLevel: "info",

		DatabasePath:   DefaultDatabasePath,
		VideosPath:     DefaultVidsPath,
		ImagesPath:     DefaultImagesPath,
		ThumbnailsPath: DefaultThumbnailsPath,

		WebServerConfig: webserverConfig,
	}
}

func setupDatabase(config Config) (*gorm.DB, error) {
	db, err := database.Open(config.DatabasePath)
	if err != nil {
		return nil, err
	}

	cpuCount := runtime.NumCPU() - 1
	if cpuCount <= 0 {
		cpuCount = 1
	}

	log.WithField("worker-count", cpuCount).
		Info("Starting FFMPeg workers")

	requests := make(chan database.FFMPegRequest)
	go database.StartFFMPegWorkerPools(cpuCount, requests)
	go database.ScanForNewFileVideos(db, config.VideosPath, config.ThumbnailsPath, requests)
	go database.ScanForNewFileImages(db, config.ImagesPath, config.ThumbnailsPath, requests)
	return db, nil
}

func setLogLevel(levelName string) {
	switch levelName {
	case "panic":
		log.SetLevel(log.PanicLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "trace":
		log.SetLevel(log.TraceLevel)
	default:
		log.SetLevel(log.InfoLevel)
		log.WithField("level", levelName).
			Warn("Unknown log level name")
		levelName = "info"
	}

	log.WithField("level", levelName).
		Info("Using log level")
}

func fileConfig(filePath string, config Config) Config {
	configFile, err := ini.Load(filePath)
	if err != nil {
		return config
	}

	mediaSection := configFile.Section("media")
	return Config{
		DatabasePath:    mediaSection.Key("database").MustString(config.DatabasePath),
		ImagesPath:      mediaSection.Key("images").MustString(config.ImagesPath),
		VideosPath:      mediaSection.Key("videos").MustString(config.VideosPath),
		ThumbnailsPath:  mediaSection.Key("thumbnails").MustString(config.ThumbnailsPath),
		WebServerConfig: webserver.FileWebServerConfig(configFile, config.WebServerConfig),
	}
}

func commandLineConfig(config Config) Config {
	logLevel := flag.String("log-level", config.LogLevel,
		"Log level (panic, fatal, error, warn, info, debug and trace)")

	databasePath := flag.String("database", config.DatabasePath, "Path to sqlite database file")
	videosPath := flag.String("vids", config.VideosPath, "Path to videos")
	imagesPath := flag.String("images", config.ImagesPath, "Path to images")
	thumbnailsPath := flag.String("thumbnails", config.ThumbnailsPath, "Path to thumbnails")

	webserverConfig := webserver.CommandLineWebServerConfig(config.WebServerConfig)
	return Config{
		LogLevel: *logLevel,

		DatabasePath:   *databasePath,
		VideosPath:     *videosPath,
		ImagesPath:     *imagesPath,
		ThumbnailsPath: *thumbnailsPath,

		WebServerConfig: webserverConfig,
	}
}

func main() {
	config := defaultConfig()
	config = fileConfig("/etc/watch-party.conf", config)
	config = commandLineConfig(config)
	setLogLevel(config.LogLevel)

	db, err := setupDatabase(config)
	if err != nil {
		panic(err)
	}

	clients := make(chan Client)
	serverMessages := make(chan ServerMessage)

	go StartSocketServer(config.WebServerConfig, clients)
	go StartServer(db, serverMessages)
	ListenForNewClients(clients, serverMessages)
}
