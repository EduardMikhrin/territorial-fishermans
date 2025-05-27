package service

import (
	"context"
	"github.com/EduardMikhrin/territorial-fishermans/cmd/utils"
	"github.com/EduardMikhrin/territorial-fishermans/internal/api"
	types "github.com/EduardMikhrin/territorial-fishermans/internal/api/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os/signal"
	"syscall"
)

const syncFlag = "sync"

var syncEnabled bool

var RunListenersCmd = &cobra.Command{
	Use:   "all",
	Short: "Run the relayer service with all listeners",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := utils.ConfigFromFlags(cmd)
		if err != nil {
			return errors.Wrap(err, "failed to get config from flags")
		}
		ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
		defer cancel()

		if err := api.RunServer(ctx, cfg); err != nil {
			if errors.Is(err, types.ErrStopped) {
				cfg.Log().Info("service stopped")
				return nil
			}
			return errors.Wrap(err, "failed to start server")
		}

		return nil
	},
}

func registerSyncFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&syncEnabled, syncFlag, "s", syncEnabled, "Sync enabled/disabled (disabled default)")
}
