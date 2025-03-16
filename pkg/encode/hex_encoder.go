package encode

import "encoding/hex"

type HexEncoder struct{}

func (h *HexEncoder) Encode(dst []byte, src []byte) {
	hex.Encode(dst, src)
}

func (h *HexEncoder) EncodeToString(src []byte) string {
	return hex.EncodeToString(src)
}
