package main

import (
	"get.porter.sh/mixin/tofu/pkg/tofu"
	"github.com/spf13/cobra"
)

func buildSchemaCommand(m *tofu.Mixin) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schema",
		Short: "Print the json schema for the mixin",
		Run: func(cmd *cobra.Command, args []string) {
			m.PrintSchema()
		},
	}
	return cmd
}
