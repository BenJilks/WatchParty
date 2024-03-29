package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

type MonkeyActionMessage struct {
	Action string `json:"action"`
	Token  string `json:"token"`
}

type ChatMessage struct {
	Message string `json:"message"`
}

type RequestPlayMessage struct {
	Playing   bool    `json:"playing"`
	Progress  float64 `json:"progress"`
	VideoFile *string `json:"video"`
}

type RequestImageMessage struct {
	File string `json:"file"`
	Name string `json:"name"`
}

func handleClient(client Client, serverMessage chan<- ServerMessage) {
	serverMessage <- ServerMessage{
		Type:   ServerMessageJoin,
		Client: &client,
	}

	for message := range client.Messages {
		log.WithFields(log.Fields{
			"token": *client.Token,
			"type":  message.Type,
		}).Trace("Got message")

		switch message.Type {
		case MessageMonkeyAction:
			var clapMessage MonkeyActionMessage
			_ = json.Unmarshal(message.Data, &clapMessage)

			serverMessage <- ServerMessage{
				Type:   ServerMessageMonkeyAction,
				Token:  client.Token,
				Action: clapMessage.Action,
			}

		case MessageChat:
			var chatMessage ChatMessage
			_ = json.Unmarshal(message.Data, &chatMessage)

			serverMessage <- ServerMessage{
				Type:    ServerMessageChat,
				Token:   client.Token,
				Message: chatMessage.Message,
			}

		case MessageVideoList:
			serverMessage <- ServerMessage{
				Type:  ServerMessageVideoList,
				Token: client.Token,
			}

		case MessageImageList:
			serverMessage <- ServerMessage{
				Type:  ServerMessageImageList,
				Token: client.Token,
			}

		case MessageRequestPlay:
			var requestMessage RequestPlayMessage
			_ = json.Unmarshal(message.Data, &requestMessage)

			serverMessage <- ServerMessage{
				Type:     ServerMessageRequestPlay,
				Token:    client.Token,
				Playing:  requestMessage.Playing,
				Progress: requestMessage.Progress,
				File:     requestMessage.VideoFile,
			}

		case MessageRequestImage:
			var requestMessage RequestImageMessage
			_ = json.Unmarshal(message.Data, &requestMessage)

			serverMessage <- ServerMessage{
				Type:  ServerMessageRequestImage,
				Token: client.Token,
				File:  &requestMessage.File,
				Name:  requestMessage.Name,
			}

		case MessageReady:
			serverMessage <- ServerMessage{
				Type:  ServerMessageReady,
				Token: client.Token,
			}

		case MessageDisconnect:
			if client.Token != nil {
				serverMessage <- ServerMessage{
					Type:  ServerMessageLeave,
					Token: client.Token,
				}
			}
			return

		default:
			log.Warnf("Uknown message type '%s'", message.Type)
		}
	}
}

func ListenForNewClients(clients <-chan Client, serverMessage chan<- ServerMessage) {
	for client := range clients {
		go handleClient(client, serverMessage)
	}
}
