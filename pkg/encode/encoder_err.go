package encode

import "fmt"

type Error struct {
	Type    string
	Message string
}

func NewError[T any](errMsg string) *Error {
	return &Error{
		Type:    fmt.Sprintf("%T", *new(T)),
		Message: errMsg,
	}
}

func (e *Error) Error() string {
	if e.Type == "" {
		return e.Message
	}

	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

var (
	ErrInvalidExtension = NewError[FileEncoder]("invalid extension")
	ErrInvalidFilePath  = NewError[FileEncoder]("invalid file path")
	ErrInvalidFile      = NewError[FileEncoder]("invalid file")
	ErrFilePathNotSet   = NewError[FileEncoder]("file path not set")
	ErrSourceTextNotSet = NewError[TextEncoder]("source text not set")
	ErrEncoderNotSet    = NewError[any]("encoder not set")
)
