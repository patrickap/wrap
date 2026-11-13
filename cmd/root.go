package cmd

import (
	"fmt"
	"strings"

	"github.com/patrickap/wrap/internal/builder"
	"github.com/patrickap/wrap/internal/config"
	"github.com/patrickap/wrap/internal/executor"
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	var configFlag string
	var verboseFlag bool
	var debugFlag bool

	root := &cobra.Command{
		Use:           "wrap [-c wrap.config.yaml] [-v | -d] <program> [arguments...]",
		Short:         "Run a CLI with configured default arguments",
		Args:          cobra.MinimumNArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(command *cobra.Command, args []string) error {
			currentConfig, err := config.Load(configFlag)

			if err != nil {
				return err
			}

			program := args[0]
			originalArgs := args[1:]
			resolvedArgs := builder.Build(program, originalArgs, currentConfig)
			resolvedCommand := strings.Join(append([]string{program}, resolvedArgs...), " ")

			if verboseFlag {
				_, err := fmt.Fprintln(command.ErrOrStderr(), resolvedCommand)

				if err != nil {
					return err
				}
			}

			if debugFlag {
				_, err := fmt.Fprintln(command.OutOrStdout(), resolvedCommand)
				return err
			}

			return executor.RunStdio(program, resolvedArgs)
		},
	}

	root.Flags().StringVarP(&configFlag, "config", "c", "", "sets configuration file path")
	root.Flags().BoolVarP(&verboseFlag, "verbose", "v", false, "outputs more verbose logs")
	root.Flags().BoolVarP(&debugFlag, "debug", "d", false, "prints resolved command without executing")
	root.Flags().SetInterspersed(false)

	return root
}

func Execute() error {
	return NewRootCommand().Execute()
}
