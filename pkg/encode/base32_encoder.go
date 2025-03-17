package encode

import "encoding/base32"

type Base32Encoder struct {
	enc *base32.Encoding
}

func NewBase32Encoder(src string) FormatEncoder {
	var encoder *base32.Encoding = base32.StdEncoding
	if len(src) == 32 {
		encoder = base32.NewEncoding(src)
	}

	return &Base32Encoder{enc: encoder}
}

func (b *Base32Encoder) Encode(dst []byte, src []byte) {
	b.enc.Encode(dst, src)
}

func (b *Base32Encoder) Decode(dst []byte, src []byte) {
	b.enc.Decode(dst, src)
}

func (b *Base32Encoder) EncodeToString(src []byte) string {
	return b.enc.EncodeToString(src)
}
