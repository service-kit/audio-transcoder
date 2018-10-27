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
	"errors"
	"log"
	"unsafe"
)

const (
	LAME_CHANNEL_MONO = 3

	LAME_VBR_OFF           = VBRMode(0)
	LAME_VBR_MT            = VBRMode(1)
	LAME_VBR_RH            = VBRMode(2)
	LAME_VBR_ABR           = VBRMode(3)
	LAME_VBR_MTRH          = VBRMode(4)
	LAME_VBR_MAX_INDICATOR = VBRMode(5)
	LAME_VBR_DEFAULT       = LAME_VBR_MTRH
)

type VBRMode int

func BytesToCUcharPoint(bytes []byte) *C.uchar {
	if nil == bytes || 0 == len(bytes) {
		return nil
	}
	return (*C.uchar)(unsafe.Pointer(&bytes[0]))
}

func BytesToCShortPoint(bytes []byte) *C.short {
	if nil == bytes || 0 == len(bytes) {
		return nil
	}
	return (*C.short)(unsafe.Pointer(&bytes[0]))
}

type Lame struct {
	lamePointer *C.lame_global_flags
}

func (l *Lame) Init(rate, channels, kbps int) {
	l.lameInit()
	l.setInSamplerate(rate)
	l.setVBRMeanBitrateKbps(kbps)
	l.setVBR(LAME_VBR_DEFAULT)
	l.setNumChannels(channels)
}

func (l *Lame) lameInit() {
	l.lamePointer = C.lame_init()
}

func (l *Lame) setInSamplerate(rate int) {
	C.lame_set_in_samplerate(l.lamePointer, C.int(rate))
}

func (l *Lame) setNumChannels(channels int) {
	C.lame_set_num_channels(l.lamePointer, C.int(channels))
}

func (l *Lame) setVBRMeanBitrateKbps(kbps int) {
	C.lame_set_VBR_mean_bitrate_kbps(l.lamePointer, C.int(kbps))
}

func (l *Lame) setVBR(vbrMode VBRMode) {
	C.lame_set_VBR(l.lamePointer, C.vbr_mode(vbrMode))
}

func (l *Lame) initParams() error {
	retCode := C.lame_init_params(l.lamePointer)
	if retCode < 0 {
		log.Printf("init params err %v", retCode)
		return errors.New("init params err")
	}
	return nil
}

func (l *Lame) LameEncodeFlush(out []byte) int {
	mp3Bytes := C.lame_encode_flush(l.lamePointer, BytesToCUcharPoint(out), C.int(len(out)))
	return int(mp3Bytes)
}

func (l *Lame) LameEncodeBufferInterleaved(in, out []byte) int {
	mp3Bytes := C.lame_encode_buffer_interleaved(l.lamePointer, BytesToCShortPoint(in), C.int(len(in)/4), BytesToCUcharPoint(out), C.int(len(out)))
	return int(mp3Bytes)
}

func NewLame(rate, channels, kbps int) *Lame {
	lame := new(Lame)
	lame.Init(rate, channels, kbps)
	return lame
}
