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

func main() {
	log.SetLevel(log.TraceLevel)

	certFile := flag.String("cert", "", "TLS cert file")
	keyFile := flag.String("key", "", "TLS key file")
	port := flag.Uint("port", 8080, "Port")

	videosPath := flag.String("vids", DefaultVidsPath, "Path to videos")
	imagesPath := flag.String("images", DefaultImagesPath, "Path to images")
	thumbnailsPath := flag.String("thumbnails", DefaultThumbnailsPath, "Path to thumbnails")
	databasePath := flag.String("database", DefaultDatabasePath, "Path to sqlite database file")
	flag.Parse()

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
