package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"errors"

	"github.com/service-kit/audio-transcoder/transcoder"
)

func PcmToMp3() error {
	f, e := os.OpenFile("../test.pcm", os.O_RDONLY, 0777)
	if nil != e {
		fmt.Println("test.pcm open err:", e.Error())
		return e
	}
	defer f.Close()
	buf := new(bytes.Buffer)
	ii, e := io.Copy(buf, f)
	if nil != e {
		fmt.Println("copy file data err:", e.Error())
		return e
	}
	fmt.Println("file size:", ii)
	transcoder := transcoder.NewPcmToMp3Transcoder(8000, transcoder.CHANNELS_NOT_SET, 16)
	if nil == transcoder {
		fmt.Println("NewPcmToMp3Transcoder return nil")
		return errors.New("NewPcmToMp3Transcoder return nil")
	}
	defer transcoder.Close()
	afterTranscodingData, e := transcoder.Transcode(buf.Bytes())
	if nil != e {
		fmt.Println("transcode err:", e.Error())
		return e
	}
	of, e := os.OpenFile("../test_out.mp3", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if nil != e {
		fmt.Println("open out file err:", e.Error())
		return e
	}
	defer of.Close()
	oi, e := of.Write(afterTranscodingData)
	if nil != e {
		fmt.Println("write file err:", e.Error())
		return e
	}
	fmt.Println("write out data size:", oi)
	return nil
}

func main() {
	e := PcmToMp3()
	if nil != e {
		fmt.Println("pcm to mp3 err:", e.Error())
	} else {
		fmt.Println("pcm to mp3 success")
	}
}
