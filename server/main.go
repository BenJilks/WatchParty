package main

func main() {
    address := "0.0.0.0:8080"
    clients := make(chan Client)
    serverMessages := make(chan ServerMessage)

    go StartSocketServer(address, clients)
    go StartServer(serverMessages)
    ListenForNewClients(clients, serverMessages)
}
