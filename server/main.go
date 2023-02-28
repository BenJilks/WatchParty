package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
)

func main() {
	certFile := flag.String("cert", "", "TLS cert file")
	keyFile := flag.String("key", "", "TLS key file")
	port := flag.Uint("port", 8080, "Port")
	videosPath := flag.String("vids", DefaultVidsPath, "Path to videos")
	flag.Parse()
	log.SetLevel(log.TraceLevel)

	videos, err := GenerateThumbnails(*videosPath)
	if err != nil {
		panic(err)
	}

	address := fmt.Sprintf("0.0.0.0:%d", *port)
	clients := make(chan Client)
	serverMessages := make(chan ServerMessage)

	go StartSocketServer(address, *certFile, *keyFile, clients)
	go StartServer(serverMessages)
	ListenForNewClients(clients, serverMessages, videos)
}
