package cmd

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/Zenika/marcel/api"
	"github.com/Zenika/marcel/config"
)

func init() {
	var cfg = config.New()

	var cmd = &cobra.Command{
		Use:   "api",
		Short: "Starts marcel's api server",
		Args:  cobra.NoArgs,

		PreRunE: preRunForServer(cfg),

		Run: func(_ *cobra.Command, _ []string) {
			a := api.New()
			a.Initialize()
			a.Start()
		},
	}

	var flags = cmd.Flags()

	commonAPIFlags(flags, cfg)

	if _, err := cfg.FlagUintP(flags, "port", "p", 8090, "Listening port", "api.port"); err != nil {
		panic(err)
	}

	if _, err := cfg.FlagString(flags, "basePath", "/api", "Base path", "api.basePath"); err != nil {
		panic(err)
	}

	if _, err := cfg.FlagBool(flags, "cors", false, "Enable CORS (all origins)", "api.cors"); err != nil {
		panic(err)
	}

	Marcel.AddCommand(cmd)
}

func commonAPIFlags(flags *pflag.FlagSet, cfg *config.Config) {
	if _, err := cfg.FlagString(flags, "dbFile", "marcel.db", "Database file", "api.dbFile"); err != nil {
		panic(err)
	}

	if _, err := cfg.FlagString(flags, "pluginsDir", "plugins", "Plugins directory", "api.pluginsDir"); err != nil {
		panic(err)
	}

	if _, err := cfg.FlagString(flags, "mediasDir", "medias", "Medias directory", "api.mediasDir"); err != nil {
		panic(err)
	}

	if _, err := cfg.FlagString(flags, "dataDir", "", "Data directory", "api.dataDir"); err != nil {
		panic(err)
	}

	if _, err := cfg.FlagBool(flags, "secure", true, "Enable secure cookies", "api.auth.secure"); err != nil {
		panic(err)
	}

	if _, err := cfg.FlagDuration(flags, "authExpiration", 8*time.Hour, "Authentication token expiration", "api.auth.expiration"); err != nil {
		panic(err)
	}

	if _, err := cfg.FlagDuration(flags, "refreshExpiration", 15*24*time.Hour, "Refresh token expiration", "api.auth.refreshExpiration"); err != nil {
		panic(err)
	}
}
