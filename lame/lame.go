package lame

/*
#cgo LDFLAGS: -L${SRCDIR} -lmp3lame -lm
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "lame.h"

#define INBUFSIZE 8192
#define MP3BUFSIZE 8192
int encode(char* in,int inlen,int rate, char* out)
{
	lame_global_flags* lame = NULL;
	int ret_code = 0;
	short input_buffer[INBUFSIZE*2];
	int input_samples = 0;
	unsigned char mp3_buffer[MP3BUFSIZE];
	int mp3_bytes = 0;;
	lame = lame_init();
	if (lame == NULL)
	{
		printf("lame_init failed\n");
		return 0;
	}
	lame_set_in_samplerate(lame,rate);
	lame_set_num_channels(lame, MONO);
	lame_set_VBR_mean_bitrate_kbps(lame, 16);
	lame_set_VBR(lame, vbr_off);
	ret_code = lame_init_params(lame);
	if (ret_code < 0)
	{
		printf("lame_init_params returned %d\n", ret_code);
		return 0;
	}
	int inindex = 0;
	int outlen = inlen;
	int outindex = 0;
	int readSize = 0;
	do
	{
		readSize = (inlen - inindex > 4 * INBUFSIZE ? 4 * INBUFSIZE : inlen - inindex);
		printf("read_size:%d\n",readSize);
		printf("begin copy to input_buffer\n");
		memcpy(input_buffer, in + inindex, readSize);
		printf("end copy to input_buffer\n");
		if (readSize == 0)
		{
			mp3_bytes = lame_encode_flush(lame, mp3_buffer, MP3BUFSIZE);
		}else
		{
			mp3_bytes = lame_encode_buffer_interleaved(lame, input_buffer, readSize/4, mp3_buffer, MP3BUFSIZE);
		}
		printf("mp3_bytes:%d\n",mp3_bytes);
		printf("begin copy to mp3_buffer\n");
		memcpy(out + outindex, mp3_buffer, mp3_bytes);
		printf("end copy to mp3_buffer\n");
		outindex += mp3_bytes;
		inindex += readSize;
	} while (readSize != 0);
	lame_close(lame);
	return outindex;
}
*/
import "C"
import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"unsafe"
)

func PcmToMp3(in []byte) (out []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("pcm to mp3 err:", e)
		}
	}()
	if nil == in {
		return nil, errors.New("input data is nil")
	}
	inlen := len(in)
	if inlen == 0 {
		return nil, errors.New("input data is empty")
	}
	fmt.Printf("input len:%v\n", inlen)
	//ADD 8192 bytes Avoid PCM file streams that are less than the length of the mp3 header
	buf := make([]byte, inlen+8192)
	outSize := C.encode((*C.char)(unsafe.Pointer(&in[0])), C.int(inlen), C.int(8000), (*C.char)(unsafe.Pointer(&buf[0])))
	if outSize == 0 {
		return nil, errors.New("data is null")
	}
	fmt.Printf("output len:%v\n", outSize)
	if int(outSize) > len(buf) {
		return nil, errors.New("encode result is err,outsize:" + strconv.Itoa(int(outSize)))
	}
	out = buf[0:outSize]
	return out, nil
}

type LameChannelsType int

const (
	LAME_CHANNEL_MONO = 3

	// lame bitrate control mode
	// CBR
	LAME_VBR_OFF = VBRMode(0)
	LAME_VBR_MT  = VBRMode(1)
	LAME_VBR_RH  = VBRMode(2)
	// ABR
	LAME_VBR_ABR = VBRMode(3)
	// VBR
	LAME_VBR_MTRH = VBRMode(4)
	// don't use
	LAME_VBR_MAX_INDICATOR = VBRMode(5)
	LAME_VBR_DEFAULT       = LAME_VBR_MTRH

	// lame channels type
	// use input audio channels number
	LAME_CHANNELS_NOT_SET = LameChannelsType(4)
	// use mone
	LAME_CHANNELS_MONO = LameChannelsType(3)
	// use joint stereo
	LAME_CHANNELS_JOINT_STEREO = LameChannelsType(1)
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
	lgfp *C.lame_global_flags
}

func (l *Lame) Init(rate int, channels LameChannelsType, kbps int) {
	l.lameInit()
	l.setInSamplerate(rate)
	l.setVBRMeanBitrateKbps(kbps)
	l.setVBR(LAME_VBR_DEFAULT)
	l.setNumChannels(channels)
	l.initParams()
}

func (l *Lame) lameInit() {
	l.lgfp = C.lame_init()
}

func (l *Lame) setInSamplerate(rate int) {
	C.lame_set_in_samplerate(l.lgfp, C.int(rate))
}

func (l *Lame) setNumChannels(channels LameChannelsType) {
	C.lame_set_num_channels(l.lgfp, C.int(channels))
}

func (l *Lame) setVBRMeanBitrateKbps(kbps int) {
	C.lame_set_VBR_mean_bitrate_kbps(l.lgfp, C.int(kbps))
}

func (l *Lame) setVBR(vbrMode VBRMode) {
	C.lame_set_VBR(l.lgfp, C.vbr_mode(vbrMode))
}

func (l *Lame) initParams() error {
	retCode := C.lame_init_params(l.lgfp)
	if 0 != retCode {
		log.Printf("init params err %v\n", retCode)
		return errors.New("init params err,ret code " + strconv.Itoa(int(retCode)))
	}
	return nil
}

func (l *Lame) LameEncodeFlush(out []byte) int {
	mp3Bytes := C.lame_encode_flush(l.lgfp, BytesToCUcharPoint(out), C.int(len(out)))
	return int(mp3Bytes)
}

func (l *Lame) LameEncodeBufferInterleaved(in, out []byte) int {
	mp3Bytes := C.lame_encode_buffer_interleaved(l.lgfp, BytesToCShortPoint(in), C.int(len(in)/4), BytesToCUcharPoint(out), C.int(len(out)))
	return int(mp3Bytes)
}

func (l *Lame) LameClose() error {
	retCode := C.lame_close(l.lgfp)
	if 0 != retCode {
		log.Printf("init params err %v\n", retCode)
		return errors.New("lame close err,ret code " + strconv.Itoa(int(retCode)))
	}
	return nil
}

func NewLame(rate int, channels LameChannelsType, kbps int) *Lame {
	lame := new(Lame)
	lame.Init(rate, channels, kbps)
	return lame
}
