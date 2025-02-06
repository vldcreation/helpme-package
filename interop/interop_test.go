package interop

import (
	"reflect"
	"testing"
)

func TestNewInteropRunner(t *testing.T) {
	type args struct {
		interop Interop
	}
	tests := []struct {
		name string
		args args
		want InteropRunner
	}{
		{
			name: "Must JavascriptRunner",
			args: args{
				interop: Interop{
					Language: "javascript",
					FilePath: getAbs(TEST_DATA_PATH, "example.js"),
				},
			},
			want: &JavascriptRunner{Interop{
				Language: "javascript",
				FilePath: getAbs(TEST_DATA_PATH, "example.js"),
			}},
		},
		{
			name: "Must PythonRunner",
			args: args{
				interop: Interop{
					Language: "python",
					FilePath: getAbs(TEST_DATA_PATH, "example.py"),
				},
			},
			want: &PythonRunner{Interop{
				Language: "python",
				FilePath: getAbs(TEST_DATA_PATH, "example.py"),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewInteropRunner(tt.args.interop); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInteropRunner() = %v, want %v", got, tt.want)
			}
		})
	}
}
