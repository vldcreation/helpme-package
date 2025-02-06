package configurator

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name       string
		key        string
		value      string
		defaultVal string
		want       string
	}{
		{
			name:       "EnvVarExists",
			key:        "EXISTING_VAR",
			value:      "value",
			defaultVal: "default",
			want:       "value",
		},
		{
			name:       "EnvVarDoesNotExist",
			key:        "NON_EXISTING_VAR",
			value:      "",
			defaultVal: "default",
			want:       "default",
		},
		{
			name:       "EnvVarEmpty",
			key:        "EMPTY_VAR",
			value:      "",
			defaultVal: "default",
			want:       "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != "" {
				os.Setenv(tt.key, tt.value)
				defer os.Unsetenv(tt.key)
			}

			if got := GetEnv(tt.key, tt.defaultVal); got != tt.want {
				t.Errorf("GetEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}
