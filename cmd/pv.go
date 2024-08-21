package cmd

import (
	"github.com/stenstromen/pvcmount/resource"

	"github.com/spf13/cobra"
)

var pvcCmd = &cobra.Command{
	Use:   "pvc",
	Short: "Mount Persistent Volume",
	Long:  `Mount Persistent Volume`,
	RunE:  resource.ResourceUpdate,
}

func init() {
	rootCmd.AddCommand(pvcCmd)
	pvcCmd.Flags().StringP("namespace", "n", "default", "Namespace")
	pvcCmd.Flags().StringP("pvc", "p", "", "Persistent Volume")
}
