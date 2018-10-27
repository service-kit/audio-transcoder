package lame

/*
#cgo LDFLAGS: -L${SRCDIR} -lmp3lame -lm
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "lame.h"
*/
import "C"
import (
	"bytes"
	"errors"
	"log"
	"unsafe"
)

const (
	PCM_BUF_SIZE = 8192
	MP3_BUF_SIZE = 8192

	LAME_CHANNEL_MONO = 3

	LAME_VBR_OFF = iota
	LAME_VBR_MT
	LAME_VBR_RH
	LAME_VBR_ABR
	LAME_VBR_MTRH
	LAME_VBR_MAX_INDICATOR
	LAME_VBR_DEFAULT = LAME_VBR_MTRH
)

type VBRModel int
type CChar C.char
type CShort C.short
type CInt C.int
type CLameGlobalFlags C.lame_global_flags

func BytesToCCharPoint(bytes []byte) *CChar {
	if nil == bytes || 0 == len(bytes) {
		return nil
	}
	return (*CChar)(unsafe.Pointer(&bytes[0]))
}

func BytesToCShortPoint(bytes []byte) *CShort {
	if nil == bytes || 0 == len(bytes) {
		return nil
	}
	return (*CShort)(unsafe.Pointer(&bytes[0]))
}

type Lame struct {
	lamePointer *CLameGlobalFlags
	inBuf       *bytes.Buffer
	mp3Buf      []byte
}

func (l *Lame) Init() {
	l.lamePointer = C.lame_init()
	l.inBuf = new(bytes.Buffer)
	l.mp3Buf = make([]byte, MP3_BUF_SIZE)
}

func (l *Lame) SetInSamplerate(rate int) {
	C.lame_set_in_samplerate(l.lamePointer, CInt(rate))
}

func (l *Lame) SetNumChannels(channels int) {
	C.lame_set_num_channels(l.lamePointer, CInt(channels))
}

func (l *Lame) SetVBRMeanBitrateKbps(kbps int) {
	C.lame_set_VBR_mean_bitrate_kbps(l.lamePointer, CInt(kbps))
}

func (l *Lame) SetVBR(vbrMode VBRModel) {
	C.lame_set_VBR(l.lamePointer, CInt(vbrMode))
}

func (l *Lame) InitParams() error {
	retCode := C.lame_init_params(l.lamePointer)
	if retCode < 0 {
		log.Printf("init params err %v", retCode)
		return errors.New("init params err")
	}
	return nil
}

func (l *Lame) LameEncodeFlush() []byte {
	mp3Bytes := C.lame_encode_flush(l.lamePointer, BytesToCCharPoint(l.mp3Buf), CInt(MP3_BUF_SIZE))
	return l.mp3Buf[:int(mp3Bytes)]
}

func (l *Lame) LameEncodeBufferInterleaved(in []byte) []byte {
	mp3Bytes := C.lame_encode_buffer_interleaved(l.lamePointer, BytesToCShortPoint(in), CInt(len(in)/4), BytesToCCharPoint(l.mp3Buf), MP3_BUF_SIZE)
	return l.mp3Buf[:int(mp3Bytes)]
}

func NewLame() *Lame {
	lame := new(Lame)
	lame.Init()
	return lame
}
