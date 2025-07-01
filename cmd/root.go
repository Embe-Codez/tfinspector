package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "tfinspector",
	Short: "Analyze Terraform files for runtime and provider metadata",
	Long: `tfinspector recursively scans a directory of Terraform (*.tf) files
and reports required versions and providers in a structured format.`,
	Example: `  tfinspector scan ./my-terraform-project`,
}

func Execute() error {
	return RootCmd.Execute()
}

func Run() {
	if err := Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
