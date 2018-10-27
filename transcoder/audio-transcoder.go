package transcoder

import (
	"github.com/service-kit/audio-transcoder/lame"
)

type ChannelsType int

const (
	CHANNELS_NOT_SET      = ChannelsType(lame.LAME_CHANNELS_NOT_SET)
	CHANNELS_MONO         = ChannelsType(lame.LAME_CHANNEL_MONO)
	CHANNELS_JOINT_STEREO = ChannelsType(lame.LAME_CHANNELS_JOINT_STEREO)
)

func NewPcmToMp3Transcoder(rate int, channels ChannelsType, kbps int) Transcoder {
	t := new(PcmToMp3Transcoder)
	t.Init(rate, channels, kbps)
	return t
}
