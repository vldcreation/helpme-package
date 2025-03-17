package encode

type Encoder interface {
	Encode() (string, error)
	ApplyOpt(...EncoderOpt)
}

type FormatEncoder interface {
	Encode(dst []byte, src []byte)
	EncodeToString(src []byte) string
}

type FormatDecoder interface {
	Decode(dst []byte, src string) (int, error)
}
