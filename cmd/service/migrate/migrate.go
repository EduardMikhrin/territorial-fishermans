package migrate

import (
	"github.com/EduardMikhrin/territorial-fishermans/assets"
	"github.com/EduardMikhrin/territorial-fishermans/internal/config"
	"github.com/pkg/errors"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/spf13/cobra"
)

func init() {
	registerMigrateCommands(Cmd)
}

func registerMigrateCommands(cmd *cobra.Command) {
	cmd.AddCommand(migrateUpCmd)
	cmd.AddCommand(migrateDownCmd)
}

var Cmd = &cobra.Command{
	Use:   "migrate",
	Short: "Command for database migrations",
}

func execute(cfg config.Config, direction migrate.MigrationDirection) error {
	migrationsFs := &migrate.EmbedFileSystemMigrationSource{
		FileSystem: assets.Migrations,
		Root:       "migrations",
	}

	applied, err := migrate.Exec(cfg.DB().RawDB(), "postgres", migrationsFs, direction)
	if err != nil {
		return errors.Wrap(err, "failed to apply migrations")
	}

	cfg.Log().WithField("applied", applied).Info("migrations applied")

	return nil
}
