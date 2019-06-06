package config

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// Config is the root configuration object
var Config = struct {
	Port        uint
	LogLevel    log.Level
	DBFile      string
	PluginsPath string
	Auth        auth
}{
	Port:        8090,
	LogLevel:    log.InfoLevel,
	DBFile:      "marcel.db",
	PluginsPath: "plugins",
	Auth: auth{
		Secure:            true,
		AuthExpiration:    8 * time.Hour,
		RefreshExpiration: 15 * 24 * time.Hour,
	},
}
