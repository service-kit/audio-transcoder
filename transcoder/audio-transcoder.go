package transcoder

const (
	PCM_BUF_SIZE = 8192
	MP3_BUF_SIZE = 8192
)

func NewPcmToMp3Transcoder() Transcoder {
	t := new(PcmToMp3Transcoder)
	t.Init(8000, 1, 16)
	return t
}
