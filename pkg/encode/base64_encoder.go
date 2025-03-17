package encode

import "encoding/base64"

type Base64Encoder struct {
	enc *base64.Encoding
}

func NewBase64Encoder(src string) SourceEncoder {
	var encoder *base64.Encoding = base64.StdEncoding
	if len(src) == 64 {
		encoder = base64.NewEncoding(src)
	}

	return &Base64Encoder{enc: encoder}
}

func (b *Base64Encoder) Encode(dst []byte, src []byte) {
	b.enc.Encode(dst, src)
}

func (b *Base64Encoder) Decode(dst []byte, src []byte) {
	b.enc.Decode(dst, src)
}

func (b *Base64Encoder) EncodeToString(src []byte) string {
	return b.enc.EncodeToString(src)
}
