package encode

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.design/x/clipboard"
)

var (
	FileExtensionAllowed = []string{".png", ".jpg", ".jpeg", ".pdf", ".go"}
)

type FileEncoder struct {
	fpath           string
	copyToClipboard bool
	formatEncoder   FormatEncoder
}

func NewFileEncoder(fpath string, opts ...EncoderOpt) Encoder {
	i := &FileEncoder{
		fpath: fpath,
	}
	for _, opt := range opts {
		opt(i)
	}
	return i
}

func (i *FileEncoder) ApplyOpt(opts ...EncoderOpt) {
	for _, opt := range opts {
		opt(i)
	}
}

func (i *FileEncoder) Encode() (string, error) {
	if i.fpath == "" {
		return "", ErrFilePathNotSet
	}

	if i.formatEncoder == nil {
		return "", ErrEncoderNotSet
	}
	return i.encode()
}

func (i *FileEncoder) encode() (string, error) {
	path := filepath.Clean(i.fpath)

	ext := filepath.Ext(path)
	if !strings.Contains(strings.Join(FileExtensionAllowed, ""), ext) {
		return "", fmt.Errorf("file extension not allowed, allowed: %s", strings.Join(FileExtensionAllowed, ", "))
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("image does not exist")
	}

	file, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	res := i.formatEncoder.EncodeToString(file)

	if i.copyToClipboard {
		err = i.copyFileToCliboard(res)
		if err != nil {
			return "", err
		}
	}

	return res, nil
}

func (i *FileEncoder) copyFileToCliboard(text string) error {
	err := clipboard.Init()
	if err != nil {
		return err
	}

	changed := clipboard.Write(clipboard.FmtText, []byte(text))
	select {
	case <-changed:
		return fmt.Errorf(`"text data" is no longer available from clipboard.`)
	default:
		return nil
	}
}
