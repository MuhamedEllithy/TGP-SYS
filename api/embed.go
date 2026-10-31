// Package api embeds the authoritative contract and local documentation shell.
package api

import _ "embed"

//go:embed openapi.yaml
var Contract []byte

//go:embed docs.html
var Docs []byte
