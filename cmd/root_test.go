package cmd

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestRootCommandDebugOutput(t *testing.T) {
	configPath := filepath.Join("testdata", "wrap.config.yaml")
	programFlags := "--prefix long -p short --UPPERCASE --[special] --bool-true --camelCase --float 0.99 --int 1 --int-negative -1 --int-zero 0 --nested -1 --nested 0 --nested 1 --nested 0.99 --nested  --nested one two three --nested [] --nested {} --nested --separator=equals --str-empty  --str-list [] --str-map {} --str-words one two three"
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

func TestRootCommandLogRunsCommand(t *testing.T) {
	touch, err := exec.LookPath("touch")
	if err != nil {
		t.Skip("touch command is unavailable")
	}

	targetPath := filepath.Join(t.TempDir(), "command-ran")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	root := NewRootCommand()
	root.SetArgs([]string{"-c", filepath.Join("testdata", "wrap.config.yaml"), "--verbose", touch, targetPath})
	root.SetOut(&stdout)
	root.SetErr(&stderr)

	if err := root.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if _, err := os.Stat(targetPath); err != nil {
		t.Fatalf("target command did not create its file: %v", err)
	}
	if got, want := stderr.String(), touch+" "+targetPath+"\n"; got != want {
		t.Fatalf("log output = %q, want %q", got, want)
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
}
