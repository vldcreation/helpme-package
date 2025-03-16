package encode

type Encoder interface {
	Encode(dst []byte, src []byte)
	EncodeToString(src []byte) string
}

type Decoder interface {
	Decode(dst []byte, src string) (int, error)
}
