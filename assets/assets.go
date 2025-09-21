package assets

import (
	"embed"
)

//go:embed rules.json
var RulesJSON []byte

//go:embed *
var Assets embed.FS
