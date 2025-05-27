package migrate

import (
	"github.com/EduardMikhrin/territorial-fishermans/cmd/utils"
	"github.com/pkg/errors"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/spf13/cobra"
)

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Migrate the database up",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := utils.ConfigFromFlags(cmd)
		if err != nil {
			return errors.Wrap(err, "failed to get config from flags")
		}

		return execute(cfg, migrate.Up)
	},
}
