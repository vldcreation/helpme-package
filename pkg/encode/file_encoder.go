package encode

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"golang.design/x/clipboard"
)

var (
	// Mapping of allowed MIME types to their corresponding file extensions
	AllowedMimeTypes = map[string][]string{
		"image/png":       {".png"},
		"image/jpeg":      {".jpg", ".jpeg"},
		"application/pdf": {".pdf"},
		"text/plain":      {".go"}, // Assuming .go files are treated as plain text
	}
)

type FileEncoder struct {
	fpath           string
	copyToClipboard bool
	withMimeType    bool
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

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("image does not exist")
	}

	file, err := os.Open(path)
	if err != nil {
		return "", err
	}

	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		return "", err
	}

	mimeType := http.DetectContentType(buffer)

	// Validate the detected MIME type against the whitelist
	if !isMimeTypeAllowed(mimeType) {
		return "", fmt.Errorf("detected MIME type %s is not allowed", mimeType)
	}

	res := i.formatEncoder.EncodeToString(buffer)

	if i.withMimeType {
		res = fmt.Sprintf("data:%s;base64,%s", mimeType, res)
	}

	if i.copyToClipboard {
		err = i.copyFileToCliboard(res)
		if err != nil {
			return "", err
		}
	}

	return res, nil
}

func isMimeTypeAllowed(mimeType string) bool {
	_, exists := AllowedMimeTypes[mimeType]
	return exists
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
