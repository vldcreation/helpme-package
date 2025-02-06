package interop

type UnknownRunner struct{ I Interop }

func (runner *UnknownRunner) Run() (string, error) {
	return "", nil
}
