package transcoder

import (
	"errors"
	"fmt"
	"github.com/service-kit/audio-transcoder/lame"
)

type PcmToMp3Transcoder struct {
	lame *lame.Lame
}

func (t *PcmToMp3Transcoder) Init(rate, channels, kbps int) error {
	t.lame = lame.NewLame(rate, channels, kbps)
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
	fmt.Println("inLen:", inLen)
	mp3Buf := make([]byte, inLen+MP3_BUF_SIZE)
	fmt.Println("mp3BufLen:", len(mp3Buf))
	readSize := 0
	mp3Bytes := 0
	fmt.Printf("readSize:%v inIndex:%v bufIndex:%v\n", readSize, inIndex, bufIndex)
	for {
		if inLen-inIndex > 2*PCM_BUF_SIZE {
			readSize = 2 * PCM_BUF_SIZE
		} else {
			readSize = inLen - inIndex
		}
		fmt.Printf("readSize:%v inIndex:%v bufIndex:%v\n", readSize, inIndex, bufIndex)
		if 0 == readSize {
			mp3Bytes = t.lame.LameEncodeFlush(mp3Buf[bufIndex:])
			fmt.Println("LameEncodeFlush mp3Bytes:", mp3Bytes)
			bufIndex += mp3Bytes
			break
		} else {
			mp3Bytes = t.lame.LameEncodeBufferInterleaved(in[inIndex:inIndex+readSize], mp3Buf[bufIndex:])
			fmt.Println("LameEncodeBufferInterleaved mp3Bytes:", mp3Bytes)
			bufIndex += mp3Bytes
			inIndex += readSize
		}
	}
	return mp3Buf[:bufIndex], nil
}
