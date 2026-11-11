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
	var debugFlag bool

	root := &cobra.Command{
		Use:           "wrap [flags] <program> [arguments...]",
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

			if debugFlag {
				resolvedCommand := append([]string{program}, resolvedArgs...)
				_, err := fmt.Fprintln(command.OutOrStdout(), strings.Join(resolvedCommand, " "))
				return err
			}

			return executor.RunStdio(program, resolvedArgs)
		},
	}

	root.Flags().StringVarP(&configFlag, "config", "c", "", "configuration file")
	root.Flags().BoolVarP(&debugFlag, "debug", "d", false, "print the resolved command without running it")
	root.Flags().SetInterspersed(false)

	return root
}

func Execute() error {
	return NewRootCommand().Execute()
}
