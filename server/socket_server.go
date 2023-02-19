package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "nhooyr.io/websocket"
)

type MessageType string

type Message struct {
    Type MessageType `json:"type"`
    Data []byte      `json:"data"`
}

type Client struct {
    Messages   <-chan Message
    Connection *websocket.Conn
    Context    *context.Context
    Token      *string
}

func (client *Client) Send(messageType MessageType, data interface{}) error {
    dataJson, err := json.Marshal(data)
    if err != nil {
        return err
    }

    messageJson, err := json.Marshal(Message{
        Type: messageType,
        Data: dataJson,
    })

    if err != nil {
        return err
    }

    return client.Connection.Write(
        *client.Context, websocket.MessageText, messageJson)
}

func handleSocketConnection(
    response http.ResponseWriter,
    request *http.Request,
    clients chan<- Client,
) {
    connection, err := websocket.Accept(response, request, nil)
    if err != nil {
        log.Println(err)
        return
    }

    requestContext := request.Context()
    messages := make(chan Message)
    defer connection.Close(websocket.StatusNormalClosure, "")
    defer close(messages)

    clients <- Client{
        Messages:   messages,
        Connection: connection,
        Context:    &requestContext,
    }

    for {
        _, content, err := connection.Read(requestContext)
        if err != nil {
            fmt.Println(err)
            break
        }

        var message Message
        if err := json.Unmarshal(content, &message); err != nil {
            fmt.Println(err)
            continue
        }

        messages <- message
    }

    messages <- Message{
        Type: MessageDisconnect,
    }
}

func handleConnection(response http.ResponseWriter, request *http.Request, clients chan<- Client) {
    path := request.URL.Path
    if path == "/socket" {
        handleSocketConnection(response, request, clients)
        return
    }

    filePath := fmt.Sprintf("../client/dist/%s", path)
    http.ServeFile(response, request, filePath)
}

func StartSocketServer(address string, clients chan<- Client) {
    handler := http.HandlerFunc(
        func(response http.ResponseWriter, request *http.Request) {
            handleConnection(response, request, clients)
        })

    err := http.ListenAndServe(address, handler)
    if err != nil {
        panic(err)
    }
}
