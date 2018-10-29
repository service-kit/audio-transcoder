# audio-transcoder

Audio Transcode Tool

###1.Pcm Transcode To Mp3

Examples Code
```
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
	fmt.Println("file size:", ii)
	transcoder := transcoder.NewPcmToMp3Transcoder(8000, transcoder.CHANNELS_NOT_SET, 16)
	if nil == transcoder {
		log.Println("NewPcmToMp3Transcoder return nil")
		return errors.New("NewPcmToMp3Transcoder return nil")
	}
	defer transcoder.Close()
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
```