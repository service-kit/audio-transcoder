package lame

import (
	"errors"
	"fmt"
)

type PcmToMp3Transcoder struct {
	lame *Lame
}

func (t *PcmToMp3Transcoder) Init(rate, channels, kbps int) error {
	t.lame = NewLame(rate, channels, kbps)
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
	fmt.Println("inLen:",inLen)
	mp3Buf := make([]byte, inLen+MP3_BUF_SIZE)
	fmt.Println("mp3BufLen:",len(mp3Buf))
	readSize := 0
	for {
		if inLen-inIndex > 4*PCM_BUF_SIZE {
			readSize = 4 * PCM_BUF_SIZE
		} else {
			readSize = inLen - inIndex
		}
		fmt.Printf("readSize:%v inIndex:%v bufIndex:%v\n",readSize,inIndex,bufIndex)
		if 0 == readSize {
			bufIndex += t.lame.LameEncodeFlush(mp3Buf[bufIndex:])
			break
		} else {
			bufIndex += t.lame.LameEncodeBufferInterleaved(in[inIndex:inIndex+readSize], mp3Buf[bufIndex:])
			inIndex += readSize
		}
	}
	return mp3Buf[:bufIndex], nil
}
