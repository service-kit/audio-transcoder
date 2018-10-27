package transcoder

type Transcoder interface {
	Init(rate, channels, kbps int) error
	Transcode([]byte) ([]byte, error)
}
