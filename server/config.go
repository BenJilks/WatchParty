package main

const (
	MessageUpdateState = MessageType("update-state")
	MessageClap        = MessageType("clap")
	MessageChat        = MessageType("chat")
	MessageVideo       = MessageType("video")
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
