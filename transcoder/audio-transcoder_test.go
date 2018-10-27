package transcoder

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestNewPcmMp3Transcoder(t *testing.T) {
	f, e := os.OpenFile("../test.pcm", os.O_RDONLY, 0777)
	if nil != e {
		t.Error("test.pcm open err:", e.Error())
		return
	}
	defer f.Close()
	buf := new(bytes.Buffer)
	ii, e := io.Copy(buf, f)
	if nil != e {
		t.Error("copy file data err:", e.Error())
		return
	}
	t.Log("file size:", ii)
	transcoder := NewPcmToMp3Transcoder(8000, CHANNELS_NOT_SET, 16)
	if nil == transcoder {
		t.Error("NewPcmToMp3Transcoder return nil")
		return
	}
	defer transcoder.Close()
	afterTranscodingData, e := transcoder.Transcode(buf.Bytes())
	if nil != e {
		t.Error("transcode err:", e.Error())
		return
	}
	of, e := os.OpenFile("../test_out.mp3", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if nil != e {
		t.Error("open out file err:", e.Error())
		return
	}
	defer of.Close()
	oi, e := of.Write(afterTranscodingData)
	if nil != e {
		t.Error("write file err:", e.Error())
		return
	}
	t.Log("write out data size:", oi)
}
