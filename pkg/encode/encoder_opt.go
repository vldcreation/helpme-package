package encode

type EncoderOpt func(T any)

func WithFpath(fpath string) EncoderOpt {
	return func(T any) {
		i, ok := T.(*ImageEncoder)
		if !ok {
			panic("T is not *ImageEncoder")
		}
		i.fpath = fpath
	}
}

func WithCopyToClipboard(copyToClipboard bool) EncoderOpt {
	return func(T any) {
		switch T := T.(type) {
		case *ImageEncoder:
			T.copyToClipboard = copyToClipboard
		case *TextEncoder:
			T.copyToClipboard = copyToClipboard
		}
	}
}

func WithEncoder(encoder Encoder) EncoderOpt {
	return func(T any) {
		switch T := T.(type) {
		case *ImageEncoder:
			T.encoder = encoder
		case *TextEncoder:
			T.encoder = encoder
		}
	}
}
