package main

import (
	"crypto/rand"
	"encoding/hex"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
	"watch-party/database"
)

type ServerMessageType int

const (
	ServerMessageJoin = ServerMessageType(iota)
	ServerMessageLeave
	ServerMessageVideoList
	ServerMessageRequestPlay
	ServerMessageReady
	ServerMessageClap
	ServerMessageChat
)

type ServerMessage struct {
	Type   ServerMessageType
	Client *Client
	Token  *string
	State  string

	Message string

	Playing  bool
	Progress float64

	VideoFile *string
}

type VideoPlaybackState struct {
	Playing            bool
	Progress           float64
	LastProgressUpdate time.Time
	VideoFile          string
}

type Server struct {
	connectedClients map[string]*Client
	stage            Stage
	videoState       VideoPlaybackState
	db               *gorm.DB
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

func (server *Server) broadcastExcept(
	except string,
	messageType MessageType,
	data interface{},
) {
	for player, client := range server.connectedClients {
		if player == except {
			continue
		}

		_ = client.Send(messageType, data)
	}
}

func (server *Server) join(client *Client) {
	token := server.generateNewToken()

	client.Token = &token
	server.connectedClients[token] = client
	server.stage.PlaceViewer(token)
	server.updateSeats()
	server.updateVideoState()

	for _, client := range server.connectedClients {
		client.Ready = false
		_ = client.Send(MessageRequestPlay, RequestPlayMessage{
			Playing:   server.videoState.Playing,
			Progress:  server.videoState.Progress,
			VideoFile: &server.videoState.VideoFile,
		})
	}
}

func (server *Server) leave(token string) {
	if _, has := server.connectedClients[token]; has {
		delete(server.connectedClients, token)
	}

	server.stage.RemovePlayer(token)
	server.updateSeats()
}

type VideoDataMessage struct {
	Name          string `json:"name"`
	VideoFile     string `json:"video_file"`
	ThumbnailFile string `json:"thumbnail_file"`
}

type VideoListMessage struct {
	Videos []VideoDataMessage `json:"videos"`
}

func (server *Server) videoList(token string) {
	client, exists := server.connectedClients[token]
	if !exists {
		return
	}

	var videos []database.Video
	if result := server.db.Find(&videos); result.Error != nil {
		log.WithError(result.Error).Error("Unable to query videos")
		return
	}

	var videoDataList = make([]VideoDataMessage, len(videos))
	for i, video := range videos {
		videoDataList[i] = VideoDataMessage{
			Name:          video.Title,
			VideoFile:     video.VideoFilePath,
			ThumbnailFile: video.ThumbnailPath,
		}
	}

	_ = client.Send(MessageVideoList, VideoListMessage{
		Videos: videoDataList,
	})
}

func (server *Server) requestPlay(message ServerMessage) {
	server.logRequest(message)

	videoFile := server.videoState.VideoFile
	if message.VideoFile != nil {
		videoFile = *message.VideoFile
	}

	server.videoState = VideoPlaybackState{
		Playing:            message.Playing,
		Progress:           message.Progress,
		VideoFile:          videoFile,
		LastProgressUpdate: time.Now(),
	}

	for _, client := range server.connectedClients {
		client.Ready = false
	}

	server.broadcastExcept(*message.Token, MessageRequestPlay, RequestPlayMessage{
		Playing:   message.Playing,
		Progress:  message.Progress,
		VideoFile: message.VideoFile,
	})
}

func (server *Server) ready(token string) {
	allClientsReady := true
	for clientToken, client := range server.connectedClients {
		if token == clientToken {
			log.WithField("token", clientToken).Info("Reported as ready")
			client.Ready = true
			continue
		}

		if !client.Ready {
			allClientsReady = false
			log.WithField("token", clientToken).Info("Is still buffering")
		}
	}

	if allClientsReady {
		log.Info("All clients are ready, playing request")
		server.broadcastExcept("", MessageReady, nil)
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

func (server *Server) logRequest(message ServerMessage) {
	videoFile := "None"
	if message.VideoFile != nil {
		videoFile = *message.VideoFile
	}

	log.WithFields(log.Fields{
		"token":    *message.Token,
		"playing":  message.Playing,
		"progress": message.Progress,
		"video":    videoFile,
	}).Info("Got play request")
}

func (server *Server) handleMessage(message ServerMessage) {
	switch message.Type {
	case ServerMessageJoin:
		server.join(message.Client)
	case ServerMessageLeave:
		server.leave(*message.Token)
	case ServerMessageVideoList:
		server.videoList(*message.Token)
	case ServerMessageRequestPlay:
		server.requestPlay(message)
	case ServerMessageReady:
		server.ready(*message.Token)
	case ServerMessageClap:
		server.clap(*message.Token, message.State)
	case ServerMessageChat:
		server.chat(*message.Token, message.Message)
	default:
		panic(message)
	}
}

func StartServer(db *gorm.DB, messages <-chan ServerMessage) {
	server := Server{
		connectedClients: map[string]*Client{},
		stage: Stage{
			seatsUsed: map[string]Seat{},
		},
		videoState: VideoPlaybackState{
			Playing:            false,
			Progress:           0,
			LastProgressUpdate: time.Now(),
			VideoFile:          "",
		},
		db: db,
	}

	for message := range messages {
		server.handleMessage(message)
	}
}
