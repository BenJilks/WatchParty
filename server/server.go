package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

type ServerMessageType int

const (
	ServerMessageBroadcast = ServerMessageType(iota)
	ServerMessageJoin
	ServerMessageLeave
	ServerMessageClap
	ServerMessageChat
)

type ServerMessage struct {
	Type    ServerMessageType
	Client  *Client
	Token   *string
	Sprite  string
	Message string
}

type Server struct {
	connectedClients map[string]*Client
	stage            Stage
}

func (server *Server) broadcast() {
	for _, client := range server.connectedClients {
		_ = client.Send(MessageUpdateState,
			server.stage.UpdateMessage(client.Token))
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
	}

	for message := range messages {
		server.handleMessage(message)
	}
}