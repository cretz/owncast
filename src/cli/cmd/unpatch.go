package cmd

import (
	"github.com/cretz/owncast/src/chrome"
	"github.com/cretz/owncast/src/util"
	"github.com/spf13/cobra"
)

func init() {
	unpatchCmd := &cobra.Command{
		Use:  "unpatch [path to chrome parent dir]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Find lib and unpatch
			lib, err := chrome.FindUnpatchableLib(args[0], chrome.LoadExistingRootCADERBytes())
			if err != nil {
				return err
			}
			util.LogInfo("Unpatching library from %v to %v", lib.Path(), lib.OrigPath())
			return lib.Unpatch()
		},
	}
	rootCmd.AddCommand(unpatchCmd)
}
