package main

import (
	"github.com/spf13/cobra"

	"get.porter.sh/mixin/tofu/pkg/tofu"
)

func buildInvokeCommand(mixin *tofu.Mixin) *cobra.Command {
	opts := tofu.InvokeOptions{}

	cmd := &cobra.Command{
		Use:   "invoke",
		Short: "Execute the invoke functionality of this mixin",
		RunE: func(cmd *cobra.Command, args []string) error {
			return mixin.Invoke(cmd.Context(), opts)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&opts.Action, "action", "", "Custom action name to invoke.")

	return cmd
}
