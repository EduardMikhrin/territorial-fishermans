package service

import (
	"github.com/EduardMikhrin/territorial-fishermans/cmd/utils"
	"github.com/spf13/cobra"
)

func init() {
	registerCommands(Cmd)
	utils.RegisterConfigFlag(Cmd)
	registerSyncFlag(Cmd)
}

var Cmd = &cobra.Command{
	Use:   "run",
	Short: "Command for running service",
}

func registerCommands(cmd *cobra.Command) {
	cmd.AddCommand(RunListenersCmd)
}
