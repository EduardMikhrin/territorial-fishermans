package service

import (
	"github.com/EduardMikhrin/territorial-fishermans/cmd/service/migrate"
	svcRun "github.com/EduardMikhrin/territorial-fishermans/cmd/service/run/service"
	"github.com/EduardMikhrin/territorial-fishermans/cmd/utils"
	"github.com/spf13/cobra"
)

func init() {
	registerServiceCommands(Cmd)
	utils.RegisterConfigFlag(Cmd)
}

func registerServiceCommands(cmd *cobra.Command) {
	cmd.AddCommand(migrate.Cmd)
	cmd.AddCommand(svcRun.Cmd)

}

var Cmd = &cobra.Command{
	Use:   "service",
	Short: "Command for running service operations",
}
