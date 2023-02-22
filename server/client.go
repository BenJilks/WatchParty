package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

type ClapMessage struct {
	State string `json:"state"`
	Token string `json:"token"`
}

type ChatMessage struct {
	Message string `json:"message"`
}

type RequestPlayMessage struct {
	Playing   bool    `json:"playing"`
	Progress  float64 `json:"progress"`
	VideoFile *string `json:"video"`
}

type VideoListMessage struct {
	Videos []VideoData `json:"videos"`
}

func handleClient(client Client, serverMessage chan<- ServerMessage, videos []VideoData) {
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
		case MessageClap:
			var clapMessage ClapMessage
			_ = json.Unmarshal(message.Data, &clapMessage)

			serverMessage <- ServerMessage{
				Type:  ServerMessageClap,
				Token: &clapMessage.Token,
				State: clapMessage.State,
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
			log.WithField("count", len(videos)).
				Info("Responded to video list request")

			_ = client.Send(MessageVideoList, VideoListMessage{
				Videos: videos,
			})

		case MessageRequestPlay:
			var requestMessage RequestPlayMessage
			_ = json.Unmarshal(message.Data, &requestMessage)

			serverMessage <- ServerMessage{
				Type:      ServerMessageRequestPlay,
				Token:     client.Token,
				Playing:   requestMessage.Playing,
				Progress:  requestMessage.Progress,
				VideoFile: requestMessage.VideoFile,
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
			panic(message)
		}
	}
}

func ListenForNewClients(clients <-chan Client, serverMessage chan<- ServerMessage, videos []VideoData) {
	for client := range clients {
		go handleClient(client, serverMessage, videos)
	}
}
