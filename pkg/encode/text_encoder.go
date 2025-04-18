package encode

import (
	"fmt"

	"golang.design/x/clipboard"
)

type TextEncoder struct {
	src             []byte
	copyToClipboard bool
	formatEncoder   FormatEncoder
}

func NewTextEncoder(text string, opts ...EncoderOpt) Encoder {
	te := &TextEncoder{
		src: []byte(text),
	}

	for _, opt := range opts {
		opt(te)
	}

	return te
}

func (t *TextEncoder) ApplyOpt(opts ...EncoderOpt) {
	for _, opt := range opts {
		opt(t)
	}
}

func (t *TextEncoder) Encode() (string, error) {
	if t.src == nil {
		return "", ErrSourceTextNotSet
	}

	if t.formatEncoder == nil {
		return "", ErrEncoderNotSet
	}

	return t.encode(), nil
}

func (t *TextEncoder) encode() string {
	encoded := t.formatEncoder.EncodeToString(t.src)
	if t.copyToClipboard {
		copyTextToClipboard(encoded)
	}
	return encoded
}

func copyTextToClipboard(text string) error {
	err := clipboard.Init()
	if err != nil {
		return err
	}

	changed := clipboard.Write(clipboard.FmtText, []byte(text))
	select {
	case <-changed:
		fmt.Println(`"text data" is no longer available from clipboard.`)
	}

	return nil
}
