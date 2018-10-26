package main

import (
	"bytes"
	"github.com/service-kit/audio-transcoder/pcm-transcoder"
	"io"
	"log"
	"os"
)

func PcmToMp3() {
	f, e := os.OpenFile("../test.pcm", os.O_RDONLY, 0777)
	if nil != e {
		log.Fatalln("test.pcm open err:", e.Error())
		return
	}
	defer f.Close()
	buf := new(bytes.Buffer)
	ii, e := io.Copy(buf, f)
	if nil != e {
		log.Fatalln("copy file data err:", e.Error())
		return
	}
	log.Fatalln("file size:", ii)
	transcoder := new(pcm_transcoder.PcmMp3Transcoder)
	afterTranscodingData, e := transcoder.Transcode(buf.Bytes())
	if nil != e {
		log.Fatalln("transcode err:", e.Error())
		return
	}
	of, e := os.OpenFile("../test_out.mp3", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if nil != e {
		log.Fatalln("open out file err:", e.Error())
		return
	}
	defer of.Close()
	oi, e := of.Write(afterTranscodingData)
	if nil != e {
		log.Fatalln("write file err:", e.Error())
		return
	}
	log.Fatalln("write out data size:", oi)
}

func main() {
	PcmToMp3()
}
