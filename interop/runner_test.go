package interop

import (
	"testing"
)

func TestRunner(t *testing.T) {
	type args struct {
		interop Interop
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Test#JavascriptRunner",
			args: args{
				interop: Interop{
					Language: "javascript",
					FilePath: getAbs(TEST_DATA_PATH, "example.js"),
				},
			},
			want: "Hello World from Javascript",
		},
		{
			name: "Test#PythonRunner",
			args: args{
				interop: Interop{
					Language: "python",
					FilePath: getAbs(TEST_DATA_PATH, "example.py"),
				},
			},
			want: "Hello World from Python",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewInteropRunner(tt.args.interop).Run()
			t.Logf("Success Invoke from %s: %+v\n", tt.args.interop.Language, got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Runner.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Runner.Run() = %v, want %v", got, tt.want)
			}
		})
	}
}
