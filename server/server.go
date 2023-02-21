package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

type ServerMessageType int

const (
	ServerMessageBroadcast = ServerMessageType(iota)
	ServerMessageJoin
	ServerMessageLeave
	ServerMessageClap
	ServerMessageChat
	ServerMessageVideo
	ServerMessageVideoChange
)

type ServerMessage struct {
	Type   ServerMessageType
	Client *Client
	Token  *string
	Sprite string

	Message string

	Playing  bool
	Progress float64

	VideoFile string
}

type Server struct {
	connectedClients map[string]*Client
	stage            Stage
	videoState       VideoState
}

func (server *Server) updateVideoState() {
	if !server.videoState.Playing {
		return
	}

	timeGone := time.Since(server.videoState.LastProgressUpdate)
	server.videoState.Progress += timeGone.Seconds()
	server.videoState.LastProgressUpdate = time.Now()
}

func (server *Server) broadcast() {
	server.updateVideoState()
	videoMessage := VideoMessage{
		Playing:  server.videoState.Playing,
		Progress: server.videoState.Progress,
	}
	videoChangeMessage := VideoChangeMessage{
		VideoFile: server.videoState.VideoFile,
	}

	for _, client := range server.connectedClients {
		_ = client.Send(MessageUpdateState, server.stage.UpdateMessage(client.Token))
		_ = client.Send(MessageVideoChange, videoChangeMessage)
		_ = client.Send(MessageVideo, videoMessage)
	}
}

func (server *Server) generateNewToken() string {
	for {
		tokenBytes := make([]byte, 16)
		if _, err := rand.Read(tokenBytes); err != nil {
			panic(err)
		}

		token := hex.EncodeToString(tokenBytes)
		if _, has := server.connectedClients[token]; has {
			continue
		}

		return token
	}
}

func (server *Server) join(client *Client) {
	token := server.generateNewToken()
	fmt.Printf("Joining watch party with token '%s'\n", token)

	client.Token = &token
	server.connectedClients[token] = client
	server.stage.PlaceViewer(token)
	server.broadcast()
}

func (server *Server) leave(token string) {
	fmt.Printf("'%s', leaving watch party\n", token)

	if _, has := server.connectedClients[token]; has {
		delete(server.connectedClients, token)
	}

	server.stage.RemovePlayer(token)
	server.broadcast()
}

func (server *Server) broadcastExcept(
	except string, messageType MessageType, data interface{},
) {
	for player, client := range server.connectedClients {
		if player == except {
			continue
		}
		_ = client.Send(messageType, data)
	}
}

type ClapResponseMessage struct {
	Sprite string `json:"sprite"`
	Row    int    `json:"row"`
	Column int    `json:"column"`
}

func (server *Server) clap(token string, sprite string) {
	seat := server.stage.SeatForPlayer(token)
	if seat == nil {
		return
	}

	server.broadcastExcept(token, MessageClap, ClapResponseMessage{
		Sprite: sprite,
		Row:    seat.Row,
		Column: seat.Column,
	})
}

type ChatResponseMessage struct {
	Message string `json:"message"`
	Row     int    `json:"row"`
	Column  int    `json:"column"`
}

func (server *Server) chat(token string, message string) {
	seat := server.stage.SeatForPlayer(token)
	if seat == nil {
		return
	}

	server.broadcastExcept("", MessageChat, ChatResponseMessage{
		Message: message,
		Row:     seat.Row,
		Column:  seat.Column,
	})
}

func (server *Server) video(token string, playing bool, progress float64) {
	videoFile := server.videoState.VideoFile
	server.videoState = VideoState{
		Playing:            playing,
		Progress:           progress,
		LastProgressUpdate: time.Now(),
		VideoFile:          videoFile,
	}

	server.broadcastExcept(token, MessageVideo, VideoMessage{
		Playing:  playing,
		Progress: progress,
	})
}

func (server *Server) videoChange(token string, videoFile string) {
	server.videoState = VideoState{
		Playing:            false,
		Progress:           0,
		LastProgressUpdate: time.Now(),
		VideoFile:          videoFile,
	}

	server.broadcastExcept(token, MessageVideoChange, VideoChangeMessage{
		VideoFile: videoFile,
	})
}

func (server *Server) handleMessage(message ServerMessage) {
	switch message.Type {
	case ServerMessageJoin:
		server.join(message.Client)
	case ServerMessageLeave:
		server.leave(*message.Token)
	case ServerMessageBroadcast:
		server.broadcast()
	case ServerMessageClap:
		server.clap(*message.Token, message.Sprite)
	case ServerMessageChat:
		server.chat(*message.Token, message.Message)
	case ServerMessageVideo:
		server.video(*message.Token, message.Playing, message.Progress)
	case ServerMessageVideoChange:
		server.videoChange(*message.Token, message.VideoFile)
	default:
		panic(message)
	}
}

func StartServer(messages <-chan ServerMessage) {
	server := Server{
		connectedClients: map[string]*Client{},
		stage: Stage{
			seatsUsed: map[string]Seat{},
		},
		videoState: VideoState{
			Playing:            false,
			Progress:           0,
			LastProgressUpdate: time.Now(),
			VideoFile:          DefaultVideoFile,
		},
	}

	for message := range messages {
		server.handleMessage(message)
	}
}
