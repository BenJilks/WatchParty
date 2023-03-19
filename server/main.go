package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"watch-party/database"
)

func setupDatabase(
	databasePath string,
	videosPath string,
	imagesPath string,
	thumbnailsPath string,
) (*gorm.DB, error) {
	db, err := database.Open(databasePath)
	if err != nil {
		return nil, err
	}

	go database.ScanForNewFileVideos(db, videosPath, thumbnailsPath)
	go database.ScanForNewFileImages(db, imagesPath, thumbnailsPath)
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

func main() {
	certFile := flag.String("cert", "", "TLS cert file")
	keyFile := flag.String("key", "", "TLS key file")
	port := flag.Uint("port", 8080, "Port")
	logLevel := flag.String("log-level", "info",
		"Log level (panic, fatal, error, warn, info, debug and trace)")

	videosPath := flag.String("vids", DefaultVidsPath, "Path to videos")
	imagesPath := flag.String("images", DefaultImagesPath, "Path to images")
	thumbnailsPath := flag.String("thumbnails", DefaultThumbnailsPath, "Path to thumbnails")
	databasePath := flag.String("database", DefaultDatabasePath, "Path to sqlite database file")

	flag.Parse()
	setLogLevel(*logLevel)

	db, err := setupDatabase(*databasePath, *videosPath, *imagesPath, *thumbnailsPath)
	if err != nil {
		panic(err)
	}

	address := fmt.Sprintf("0.0.0.0:%d", *port)
	clients := make(chan Client)
	serverMessages := make(chan ServerMessage)

	go StartSocketServer(address, *certFile, *keyFile, clients)
	go StartServer(db, serverMessages)
	ListenForNewClients(clients, serverMessages)
}
