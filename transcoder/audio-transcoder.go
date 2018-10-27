package transcoder

import (
	"github.com/service-kit/audio-transcoder/lame"
)

func NewPcmToMp3Transcoder() Transcoder {
	t := new(lame.PcmToMp3Transcoder)
	t.Init(8000, 1, 16)
	return t
}
