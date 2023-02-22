package main

const DefaultVideoFile = "[Rhythm Heaven] - Fan Club (Perfect) (English)-DNbvktlB0gU.mp4"

const (
	MessageUpdateState = MessageType("update-state")
	MessageClap        = MessageType("clap")
	MessageChat        = MessageType("chat")
	MessageVideoList   = MessageType("video-list")
	MessageRequestPlay = MessageType("request-play")
	MessageReady       = MessageType("ready")
	MessageDisconnect  = MessageType("disconnect")
)

var RowSeatCount = []int{
	16,
	16,
	14,
	12,
	10,
	8,
	6,
}
