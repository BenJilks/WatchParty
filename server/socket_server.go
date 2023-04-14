package main

import (
	"context"
	"encoding/json"
	"fmt"
	webserver "github.com/benjilks/tinywebserver"
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

func connectionHandler(clients chan<- Client, webHandler http.HandlerFunc) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		url := request.URL.Path
		log.WithField("url", url).Trace("Got HTTP request")

		if request.Header.Get("Upgrade") == "websocket" && url == "/socket" {
			handleSocketConnection(response, request, clients)
			return
		}

		if url == "/" {
			url = "/index.html"
		}

		webHandler(response, request)
	}
}

func StartSocketServer(config webserver.WebServerConfig, clients chan<- Client) {
	useTLS := config.CertFilePath != "" && config.KeyFilePath != ""
	addressFormat := "http://%s:%d"
	if useTLS {
		addressFormat = "https://%s:%d"
	}

	log.WithFields(log.Fields{
		"address": fmt.Sprintf(addressFormat, config.Address, config.Port),
		"TLS":     useTLS,
	}).Info("Started listening for connections")

	webHandler := webserver.WebHandler(config)
	handler := connectionHandler(clients, webHandler)

	bindAddress := fmt.Sprintf("%s:%d", config.Address, config.Port)
	if useTLS {
		_ = http.ListenAndServeTLS(bindAddress, config.CertFilePath, config.KeyFilePath, handler)
	} else {
		_ = http.ListenAndServe(bindAddress, handler)
	}
}
