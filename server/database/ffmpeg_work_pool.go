package database

import (
	log "github.com/sirupsen/logrus"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"path"
)

type FFMPegRequest struct {
	inputPath  string
	outputPath string

	stream   *ffmpeg.Stream
	typeName string
}

func fulfillFFMPegRequest(request FFMPegRequest) {
	name := path.Base(request.outputPath)
	if err := request.stream.Run(); err != nil {
		log.WithError(err).
			WithFields(log.Fields{"file": name, "type": request.typeName}).
			Warnf("Unable to process '%s'\n", request.inputPath)
	} else {
		log.WithFields(log.Fields{"file": name, "type": request.typeName}).
			Info("Finished processing file")
	}
}

func ffmpegWorker(requests <-chan FFMPegRequest) {
	for request := range requests {
		fulfillFFMPegRequest(request)
	}
}

func StartFFMPegWorkerPools(count int, requests <-chan FFMPegRequest) {
	for i := 0; i < count; i++ {
		go ffmpegWorker(requests)
	}
}
