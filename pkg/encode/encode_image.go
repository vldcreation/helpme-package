package encode

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.design/x/clipboard"
)

var (
	ImageExtensionAllowed = []string{".png", ".jpg", ".jpeg"}
)

func encodeImage(path string, encoder Encoder, copyToClipboard bool) (string, error) {
	path = filepath.Clean(path)

	ext := filepath.Ext(path)
	if !strings.Contains(strings.Join(ImageExtensionAllowed, ""), ext) {
		return "", fmt.Errorf("image extension not allowed")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("image does not exist")
	}

	file, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	res := encoder.EncodeToString(file)

	if copyToClipboard {
		err = copyImageToCliboard(res)
		if err != nil {
			return "", err
		}
	}

	return path, nil
}

func copyImageToCliboard(text string) error {
	err := clipboard.Init()
	if err != nil {
		return err
	}

	changed := clipboard.Write(clipboard.FmtText, []byte(text))

	select {
	case <-changed:
		println(`"text data" is no longer available from clipboard.`)
	}

	return nil
}
