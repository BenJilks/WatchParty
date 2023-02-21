package main

import "encoding/json"

type ClapMessage struct {
	Sprite string `json:"sprite"`
	Token  string `json:"token"`
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
		case MessageUpdateState:
			serverMessage <- ServerMessage{
				Type:   ServerMessageBroadcast,
				Client: &client,
				Token:  client.Token,
			}

		case MessageClap:
			var clapMessage ClapMessage
			_ = json.Unmarshal(message.Data, &clapMessage)

			serverMessage <- ServerMessage{
				Type:   ServerMessageClap,
				Token:  &clapMessage.Token,
				Sprite: clapMessage.Sprite,
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
