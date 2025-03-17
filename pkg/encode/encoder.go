package encode

type Encoder interface {
	Encode() (string, error)
	ApplyOpt(...EncoderOpt)
}

type SourceEncoder interface {
	Encode(dst []byte, src []byte)
	EncodeToString(src []byte) string
}

type SourceDecoder interface {
	Decode(dst []byte, src string) (int, error)
}
