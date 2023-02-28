package main

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
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
	Ready      bool
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
		*client.Context, websocket.MessageBinary, messageJson)
}

func handleSocketConnection(
	response http.ResponseWriter,
	request *http.Request,
	clients chan<- Client,
) {
	log.Info("Opening new web socket connection")

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
		Ready:      false,
	}

	for {
		_, content, err := connection.Read(requestContext)
		if err != nil {
			log.WithError(err).Info()
			break
		}

		var message Message
		if err := json.Unmarshal(content, &message); err != nil {
			log.WithError(err).Error("Unable to decode JSON message")
			continue
		}

		messages <- message
	}

	log.Info("Socket connection closed")
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

func StartSocketServer(address string, certFile string, keyFile string, clients chan<- Client) {
	useTLS := certFile != "" && keyFile != ""
	addressFormat := "http://%s"
	if useTLS {
		addressFormat = "https://%s"
	}

	log.WithFields(log.Fields{
		"address": fmt.Sprintf(addressFormat, address),
		"TLS":     useTLS,
	}).Info("Started listening for connections")

	handler := http.HandlerFunc(
		func(response http.ResponseWriter, request *http.Request) {
			handleConnection(response, request, clients)
		})

	if useTLS {
		_ = http.ListenAndServeTLS(address, certFile, keyFile, handler)
	} else {
		_ = http.ListenAndServe(address, handler)
	}
}
