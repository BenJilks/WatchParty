package main

import (
	"flag"
	"fmt"
)

func main() {
    certFile := flag.String("cert", "", "TLS cert file")
    keyFile := flag.String("key", "", "TLS key file")
    port := flag.Uint("port", 8080, "Port")
    flag.Parse()

    address := fmt.Sprintf("0.0.0.0:%d", *port)
    clients := make(chan Client)
    serverMessages := make(chan ServerMessage)

    go StartSocketServer(address, *certFile, *keyFile, clients)
    go StartServer(serverMessages)
    ListenForNewClients(clients, serverMessages)
}

