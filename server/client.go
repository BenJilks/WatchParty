package main

func handleClient(client Client, serverMessage chan<- ServerMessage) {
    serverMessage <- ServerMessage{
        Type:   ServerMessageJoin,
        Client: &client,
    }

    for message := range client.Messages {
        switch message.Type {
        case MessageUpdateState:
            serverMessage <- ServerMessage{
                Type:   ServerMessageBroadcast,
                Client: &client,
            }

        case MessageDisconnect:
            if client.Token != nil {
                serverMessage <- ServerMessage{
                    Type:  ServerMessageLeave,
                    Token: *client.Token,
                }
            }
            return

        default:
            panic(message)
        }
    }
}

func ListenForNewClients(clients <-chan Client, serverMessage chan<- ServerMessage) {
    for client := range clients {
        go handleClient(client, serverMessage)
    }
}
