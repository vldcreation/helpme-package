package encode

type EncoderOpt func(T any)

func WithFpath(fpath string) EncoderOpt {
	return func(T any) {
		i, ok := T.(*FileEncoder)
		if !ok {
			return
		}
		i.fpath = fpath
	}
}

func WithCopyToClipboard(copyToClipboard bool) EncoderOpt {
	return func(T any) {
		switch T := T.(type) {
		case *FileEncoder:
			T.copyToClipboard = copyToClipboard
		case *TextEncoder:
			T.copyToClipboard = copyToClipboard
		}
	}
}

func WithFormatEncoder(encoder FormatEncoder) EncoderOpt {
	return func(T any) {
		switch T := T.(type) {
		case *FileEncoder:
			T.formatEncoder = encoder
		case *TextEncoder:
			T.formatEncoder = encoder
		}
	}
}

func WithMimeType(mimeType bool) EncoderOpt {
	return func(T any) {
		i, ok := T.(*FileEncoder)
		if !ok {
			return
		}
		i.withMimeType = mimeType
	}
}
