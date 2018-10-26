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
	"strconv"
	"unsafe"
)

type PcmMp3Transcoder struct {
}

func (t *PcmMp3Transcoder) Transcode(in []byte) (out []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("pcm to mp3 err:", e)
		}
	}()
	if nil == in || 0 == len(in) {
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
