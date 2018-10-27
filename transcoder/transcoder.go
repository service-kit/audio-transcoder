package transcoder

type Transcoder interface {
	Init(rate int, channels ChannelsType, kbps int) error
	Transcode([]byte) ([]byte, error)
	Close() error
}
