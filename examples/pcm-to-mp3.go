package main

import (
	"bytes"
	"github.com/service-kit/audio-transcoder/pcm-transcoder"
	"io"
	"log"
	"os"
)

func PcmToMp3() error {
	f, e := os.OpenFile("../test.pcm", os.O_RDONLY, 0777)
	if nil != e {
		log.Println("test.pcm open err:", e.Error())
		return e
	}
	defer f.Close()
	buf := new(bytes.Buffer)
	ii, e := io.Copy(buf, f)
	if nil != e {
		log.Println("copy file data err:", e.Error())
		return e
	}
	log.Fatalln("file size:", ii)
	transcoder := new(pcm_transcoder.PcmMp3Transcoder)
	afterTranscodingData, e := transcoder.Transcode(buf.Bytes())
	if nil != e {
		log.Println("transcode err:", e.Error())
		return e
	}
	of, e := os.OpenFile("../test_out.mp3", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if nil != e {
		log.Println("open out file err:", e.Error())
		return e
	}
	defer of.Close()
	oi, e := of.Write(afterTranscodingData)
	if nil != e {
		log.Println("write file err:", e.Error())
		return e
	}
	log.Println("write out data size:", oi)
	return nil
}

func main() {
	e := PcmToMp3()
	if nil != e {
		log.Println("pcm to mp3 err:", e.Error())
	} else {
		log.Println("pcm to mp3 success")
	}
}
