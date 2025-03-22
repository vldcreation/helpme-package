package trackclipboard

import "context"

type TrackChannel interface {
	Send(ctx context.Context, msg string) error
	Close() error
}
