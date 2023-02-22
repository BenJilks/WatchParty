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

type VideoMessage struct {
	Playing  bool    `json:"playing"`
	Progress float64 `json:"progress"`
}

type VideoListMessage struct {
	Videos []VideoData `json:"videos"`
}

type VideoChangeMessage struct {
	VideoFile string `json:"video_file"`
}

func handleClient(client Client, serverMessage chan<- ServerMessage, videos []VideoData) {
	serverMessage <- ServerMessage{
		Type:   ServerMessageJoin,
		Client: &client,
	}

	for message := range client.Messages {
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

		case MessageVideo:
			var videoMessage VideoMessage
			_ = json.Unmarshal(message.Data, &videoMessage)

			serverMessage <- ServerMessage{
				Type:     ServerMessageVideo,
				Token:    client.Token,
				Playing:  videoMessage.Playing,
				Progress: videoMessage.Progress,
			}

		case MessageVideoList:
			log.WithField("count", len(videos)).
				Info("Responded to video list request")

			_ = client.Send(MessageVideoList, VideoListMessage{
				Videos: videos,
			})

		case MessageVideoChange:
			var videoChangeMessage VideoChangeMessage
			_ = json.Unmarshal(message.Data, &videoChangeMessage)

			serverMessage <- ServerMessage{
				Type:      ServerMessageVideoChange,
				Token:     client.Token,
				VideoFile: videoChangeMessage.VideoFile,
			}

		case MessageReady:
			client.Ready = true
			serverMessage <- ServerMessage{
				Type: ServerMessageReady,
			}

		case MessageWaiting:
			client.Ready = false
			serverMessage <- ServerMessage{
				Type:  ServerMessageWaiting,
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
