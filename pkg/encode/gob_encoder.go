package encode

import (
	"bytes"
	"encoding/gob"
)

type GobEncoder struct {
	enc *gob.Encoder
	dec *gob.Decoder
}

func NewGobEncoder() *GobEncoder {
	b := bytes.NewBuffer(nil)
	return &GobEncoder{
		enc: gob.NewEncoder(b),
		dec: gob.NewDecoder(b),
	}
}

// WIP
func (g *GobEncoder) Encode(src []byte) ([]byte, error) {
	return src, nil
}

func (g *GobEncoder) Decode(dst, src []byte) (int, error) {
	return len(src), nil
}
