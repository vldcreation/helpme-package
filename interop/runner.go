package interop

type InteropRunner interface {
	Run() (string, error)
}
