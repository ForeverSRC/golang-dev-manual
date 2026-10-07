// Package data embeds the clause data so gdm-cli carries its input in the binary.
package data

import "embed"

// FS is the read-only data tree: manual.json, clauses/ and i18n/.
//
//go:embed manual.json clauses i18n
var FS embed.FS
