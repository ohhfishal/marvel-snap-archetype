package assets

import (
	_ "embed"
)

//go:embed rules.json
var RulesJSON []byte
