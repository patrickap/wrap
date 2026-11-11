package cmd

import (
	"bytes"
	"path/filepath"
	"testing"
)

func TestRootCommandDebugOutput(t *testing.T) {
	configPath := filepath.Join("testdata", "config.yaml")
	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "flag types and prefixes",
			args: []string{"tool"},
			want: "tool --format=wide -v --count=2 --empty= --items=one --items=2 --items --ratio=1.5 --text=two words --true-flag --zero=0\n",
		},
		{
			name: "nested flags and args suffix",
			args: []string{"tool", "deploy", "target", "--literal=value"},
			want: "tool --format=wide -v --count=2 --empty= --items=one --items=2 --items --ratio=1.5 --text=two words --true-flag --zero=0 deploy --after --force target --literal=value --suffix=last\n",
		},
		{
			name: "unmatched command passes remaining arguments through",
			args: []string{"passthrough", "unknown", "deploy", "--literal"},
			want: "passthrough --before unknown deploy --literal\n",
		},
		{
			name: "unconfigured program",
			args: []string{"unconfigured", "one", "two"},
			want: "unconfigured one two\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var output bytes.Buffer

			root := NewRootCommand()
			root.SetArgs(append([]string{"-c", configPath, "-d"}, test.args...))
			root.SetOut(&output)

			if err := root.Execute(); err != nil {
				t.Fatalf("Execute() error = %v", err)
			}

			if got := output.String(); got != test.want {
				t.Fatalf("output = %q, want %q", got, test.want)
			}
		})
	}
}
