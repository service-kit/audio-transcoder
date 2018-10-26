package audio_transcoder

import (
	"github.com/service-kit/audio-transcoder/pcm-transcoder"
	"github.com/service-kit/audio-transcoder/transcoder"
)

func NewPcmMp3Transcoder() transcoder.Transcoder {
	t := new(pcm_transcoder.PcmMp3Transcoder)
	return t
}
