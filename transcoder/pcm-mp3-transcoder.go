package transcoder

import (
	"errors"

	"github.com/service-kit/audio-transcoder/lame"
)

const (
	PCM_BUF_SIZE = 8192
	MP3_BUF_SIZE = 8192
)

type PcmToMp3Transcoder struct {
	lame *lame.Lame
}

func (t *PcmToMp3Transcoder) Init(rate int, channels ChannelsType, kbps int) error {
	t.lame = lame.NewLame(rate, lame.LameChannelsType(channels), kbps)
	return nil
}

func (t *PcmToMp3Transcoder) Transcode(in []byte) (out []byte, err error) {
	if nil == t.lame {
		return nil, errors.New("transcoder not init")
	}
	if nil == in || 0 == len(in) {
		return nil, errors.New("in data is nil or empty")
	}
	inIndex := 0
	bufIndex := 0
	inLen := len(in)
	mp3Buf := make([]byte, inLen+MP3_BUF_SIZE)
	readSize := 0
	mp3Bytes := 0
	for {
		if inLen-inIndex > 2*PCM_BUF_SIZE {
			readSize = 2 * PCM_BUF_SIZE
		} else {
			readSize = inLen - inIndex
		}
		if 0 == readSize {
			mp3Bytes = t.lame.LameEncodeFlush(mp3Buf[bufIndex:])
			bufIndex += mp3Bytes
			break
		} else {
			mp3Bytes = t.lame.LameEncodeBufferInterleaved(in[inIndex:inIndex+readSize], mp3Buf[bufIndex:])
			bufIndex += mp3Bytes
			inIndex += readSize
		}
	}
	return mp3Buf[:bufIndex], nil
}

func (t *PcmToMp3Transcoder) Close() error {
	if nil != t.lame {
		return t.lame.LameClose()
	}
	return errors.New("transcoder not init")
}
