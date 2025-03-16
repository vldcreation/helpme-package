package encode

import (
	"fmt"

	"golang.design/x/clipboard"
)

func encodeText(src []byte, encoder Encoder, copyToClipboard bool) string {
	encoded := encoder.EncodeToString(src)
	if copyToClipboard {
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
