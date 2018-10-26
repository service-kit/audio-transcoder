package audio_transcoder

import (
	"github.com/service-kit/audio-transcoder/lame"
	"github.com/service-kit/audio-transcoder/transcoder"
)

func NewPcmMp3Transcoder() transcoder.Transcoder {
	t := new(lame.PcmMp3Transcoder)
	return t
}
