package transcoder

type Transcoder interface {
	Transcode([]byte) ([]byte, error)
}
