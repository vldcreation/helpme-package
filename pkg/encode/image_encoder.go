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

type ImageEncoder struct {
	fpath           string
	copyToClipboard bool
	encoder         Encoder
}

func NewImageEncoder(fpath string, opts ...EncoderOpt) *ImageEncoder {
	i := &ImageEncoder{}
	for _, opt := range opts {
		opt(i)
	}
	return i
}

func (i *ImageEncoder) Encode() (string, error) {
	if i.fpath == "" {
		return "", ErrFilePathNotSet
	}

	if i.encoder == nil {
		return "", ErrEncoderNotSet
	}
	return i.encode()
}

func (i *ImageEncoder) encode() (string, error) {
	path := filepath.Clean(i.fpath)

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

	res := i.encoder.EncodeToString(file)

	if i.copyToClipboard {
		err = i.copyImageToCliboard(res)
		if err != nil {
			return "", err
		}
	}

	return path, nil
}

func (i *ImageEncoder) copyImageToCliboard(text string) error {
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
