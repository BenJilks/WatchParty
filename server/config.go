package main

const StaticFilesPath = "../client/dist"
const DefaultVidsPath = StaticFilesPath + "/vids"
const DefaultImagesPath = StaticFilesPath + "/images"
const DefaultThumbnailsPath = StaticFilesPath + "/thumbnails"
const DefaultDatabasePath = StaticFilesPath + "/watch-party.db"

const (
	MessageUpdateState  = MessageType("update-state")
	MessageMonkeyAction = MessageType("monkey-action")
	MessageChat         = MessageType("chat")
	MessageVideoList    = MessageType("video-list")
	MessageImageList    = MessageType("image-list")
	MessageRequestPlay  = MessageType("request-play")
	MessageRequestImage = MessageType("request-image")
	MessageReady        = MessageType("ready")
	MessageDisconnect   = MessageType("disconnect")
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
