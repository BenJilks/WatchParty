package main

import (
	"crypto/rand"
	"encoding/hex"
	log "github.com/sirupsen/logrus"
	"time"
)

type ServerMessageType int

const (
	ServerMessageJoin = ServerMessageType(iota)
	ServerMessageLeave
	ServerMessageClap
	ServerMessageChat
	ServerMessageVideo
	ServerMessageVideoChange
	ServerMessageReady
	ServerMessageWaiting
)

type ServerMessage struct {
	Type   ServerMessageType
	Client *Client
	Token  *string
	State  string

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

func (server *Server) updateSeats() {
	for _, client := range server.connectedClients {
		_ = client.Send(MessageUpdateState, server.stage.UpdateMessage(client.Token))
	}
}

func (server *Server) generateNewToken() string {
	for {
		tokenBytes := make([]byte, 8)
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

	client.Token = &token
	server.connectedClients[token] = client
	server.stage.PlaceViewer(token)
	server.updateSeats()
	server.updateVideoState()

	videoMessage := VideoMessage{
		Playing:  server.videoState.Playing,
		Progress: server.videoState.Progress,
	}
	videoChangeMessage := VideoChangeMessage{
		VideoFile: server.videoState.VideoFile,
	}
	_ = client.Send(MessageVideoChange, videoChangeMessage)
	_ = client.Send(MessageVideo, videoMessage)
}

func (server *Server) leave(token string) {
	if _, has := server.connectedClients[token]; has {
		delete(server.connectedClients, token)
	}

	server.stage.RemovePlayer(token)
	server.updateSeats()
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
	State  string `json:"state"`
	Row    int    `json:"row"`
	Column int    `json:"column"`
}

func (server *Server) clap(token string, state string) {
	seat := server.stage.SeatForPlayer(token)
	if seat == nil {
		return
	}

	server.broadcastExcept(token, MessageClap, ClapResponseMessage{
		State:  state,
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

	log.WithFields(log.Fields{
		"token":   token,
		"message": message,
		"row":     seat.Row,
		"seat":    seat.Column,
	}).Trace("Send chat message")

	server.broadcastExcept("", MessageChat, ChatResponseMessage{
		Message: message,
		Row:     seat.Row,
		Column:  seat.Column,
	})
}

func (server *Server) video(token string, playing bool, progress float64) {
	if server.videoState.Playing && !playing {
		log.WithField("token", token).Trace("Paused")
	}
	if !server.videoState.Playing && playing {
		log.WithField("token", token).Trace("Resumed")
	}
	if server.videoState.Progress != progress {
		log.WithFields(log.Fields{
			"token":    token,
			"progress": progress,
		}).Trace("Seeked")
	}

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
		Playing:            true,
		Progress:           0,
		LastProgressUpdate: time.Now(),
		VideoFile:          videoFile,
	}

	log.WithFields(log.Fields{
		"token": token,
		"video": videoFile,
	}).Info("Changed video")

	for _, client := range server.connectedClients {
		client.Ready = false
	}

	server.broadcastExcept(token, MessageVideoChange, VideoChangeMessage{
		VideoFile: videoFile,
	})
}

func (server *Server) ready() {
	for token, client := range server.connectedClients {
		if !client.Ready {
			log.WithField("token", token).Info("Is still buffering")
			return
		}
	}

	server.broadcastExcept("", MessageReady, nil)
}

func (server *Server) waiting(token string) {
	log.WithField("token", token).Info("Waiting for buffering")
	server.broadcastExcept(token, MessageSyncing, nil)
}

func (server *Server) handleMessage(message ServerMessage) {
	switch message.Type {
	case ServerMessageJoin:
		server.join(message.Client)
	case ServerMessageLeave:
		server.leave(*message.Token)
	case ServerMessageClap:
		server.clap(*message.Token, message.State)
	case ServerMessageChat:
		server.chat(*message.Token, message.Message)
	case ServerMessageVideo:
		server.video(*message.Token, message.Playing, message.Progress)
	case ServerMessageVideoChange:
		server.videoChange(*message.Token, message.VideoFile)
	case ServerMessageReady:
		server.ready()
	case ServerMessageWaiting:
		server.waiting(*message.Token)
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
