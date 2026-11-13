package cmd

import (
	"bytes"
	"path/filepath"
	"testing"
)

func TestRootCommandDebugOutput(t *testing.T) {
	configPath := filepath.Join("testdata", "wrap.config.yaml")
	programFlags := "--prefix=long -p=short --[special] --bool-true --camelcase --float=0.99 --int=1 --int-negative=-1 --int-zero=0 --nested=-1 --nested=0 --nested=1 --nested=0.99 --nested= --nested=one two three --nested=[] --nested={} --nested --str-empty= --str-list=[] --str-map={} --str-words=one two three --uppercase"
	tests := []struct {
		name string
		args []string
		want string
	}{
		// program or command with flags
		{
			name: "program handles flags",
			args: []string{"program", "--flag"},
			want: "program " + programFlags + " --flag --after-args\n",
		},
		{
			name: "program with commands handles flags",
			args: []string{"program", "command", "subcommand", "--flag"},
			want: "program " + programFlags + " command --after-command subcommand --after-subcommand --flag\n",
		},
		{
			name: "program with commands handles flags in between",
			args: []string{"program", "command", "--flag", "subcommand"},
			want: "program " + programFlags + " command --after-command --flag subcommand --after-args\n",
		},

		// program or command with arguments
		{
			name: "program handles args",
			args: []string{"program", "args"},
			want: "program " + programFlags + " args --after-args\n",
		},
		{
			name: "program with commands handles args",
			args: []string{"program", "command", "args"},
			want: "program " + programFlags + " command --after-command args --after-args\n",
		},
		{
			name: "program with commands handles args in between",
			args: []string{"program", "args", "subcommand"},
			want: "program " + programFlags + " args subcommand --after-args\n",
		},
		{
			name: "program with commands handles args and flags",
			args: []string{"program", "args", "--flag", "subcommand", "--flag"},
			want: "program " + programFlags + " args --flag subcommand --flag --after-args\n",
		},

		// unknown program or command
		{
			name: "unknown program passes through",
			args: []string{"unknown"},
			want: "unknown\n",
		},
		{
			name: "unknown program with args and flags passes through",
			args: []string{"unknown", "args", "--flag"},
			want: "unknown args --flag\n",
		},
		{
			name: "unknown command passes through",
			args: []string{"program", "command", "subcommand", "unknown"},
			want: "program " + programFlags + " command --after-command subcommand --after-subcommand unknown\n",
		},
		{
			name: "unknown command with args and flags passes through",
			args: []string{"program", "command", "subcommand", "unknown", "args", "--flag"},
			want: "program " + programFlags + " command --after-command subcommand --after-subcommand unknown args --flag\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var output bytes.Buffer

			root := NewRootCommand()
			root.SetArgs(append([]string{"-c", configPath, "-d"}, test.args...))
			root.SetOut(&output)

			err := root.Execute()

			if err != nil {
				t.Fatalf("Execute() error = %v", err)
			}

			got := output.String()

			if got != test.want {
				t.Fatalf("output = %q, want %q", got, test.want)
			}
		})
	}
}
