// Package core carries the tool-neutral files of the factory that the
// binary needs at run time.
package core

import _ "embed"

// QuickValidate is the validator of the skill authoring-skills, the one
// source for the rules of skills (ADR 0001).
//
//go:embed skills/authoring-skills/scripts/quick_validate.py
var QuickValidate []byte
