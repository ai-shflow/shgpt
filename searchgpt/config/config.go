package config

import (
	"embed"
)

//go:embed config.yml
var ConfigFile embed.FS
