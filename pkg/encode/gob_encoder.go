package encode

import (
	"bytes"
	"encoding/gob"
)

type GobEncoder struct {
	buf *bytes.Buffer
	enc *gob.Encoder
	dec *gob.Decoder
}

func NewGobEncoder() FormatEncoder {
	b := bytes.NewBuffer(nil)
	return &GobEncoder{
		buf: b,
		enc: gob.NewEncoder(b),
		dec: gob.NewDecoder(b),
	}
}

// WIP
func (g *GobEncoder) Encode(src []byte, dst []byte) {
}

func (g *GobEncoder) EncodeToString(src []byte) string {
	if err := g.enc.Encode(src); err != nil {
		return ""
	}

	return g.buf.String()
}

func (g *GobEncoder) Decode(dst, src []byte) (int, error) {
	return len(src), nil
}
