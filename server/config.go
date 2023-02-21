package main

const DefaultVideoFile = "[Rhythm Heaven] - Fan Club (Perfect) (English)-DNbvktlB0gU.mp4"

const (
	MessageUpdateState = MessageType("update-state")
	MessageClap        = MessageType("clap")
	MessageChat        = MessageType("chat")
	MessageVideo       = MessageType("video")
	MessageVideoList   = MessageType("video-list")
	MessageVideoChange = MessageType("video-change")
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
