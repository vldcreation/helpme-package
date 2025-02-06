package interop

import "path/filepath"

func getAbs(path ...string) string {
	joinPath := filepath.Join(path...)
	abs, err := filepath.Abs(joinPath)
	if err != nil {
		panic(err)
	}
	return abs
}
