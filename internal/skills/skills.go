// Package skills runs the validator of the skill authoring-skills on each
// skill under a set of folders and reports its findings.
package skills

import (
	"bufio"
	"bytes"
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/HannesT117/software-factory/core"
	"github.com/HannesT117/software-factory/internal/report"
)

// python is the interpreter that runs the validator.
const python = "python3"

// probe exits with 0 when Python is 3.10 or later and PyYAML is installed.
const probe = `import sys
try:
    import yaml
except ImportError:
    sys.exit(3)
sys.exit(0 if sys.version_info >= (3, 10) else 3)`

// boot runs the validator from the folder in FACTORY_VALIDATOR_DIR. It drops
// the current folder from the module path first, so that a file in the
// checked folder cannot replace a module.
const boot = `import os, sys
sys.path[:] = [p for p in sys.path if p not in ('', '.')]
sys.path.insert(0, os.environ['FACTORY_VALIDATOR_DIR'])
from scripts.quick_validate import main
sys.exit(main(sys.argv[1:]))`

const prefix = "quick_validate: "

// Validate runs the validator on each skill under dirs and returns the exit
// code. With strict, warnings count as errors.
func Validate(r *report.Report, dirs []string, strict bool) int {
	var found []string
	for _, dir := range dirs {
		list, err := find(dir)
		if err != nil {
			return r.Abort(dir, "unreadable", err.Error())
		}
		found = append(found, list...)
	}

	if out, err := exec.Command(python, "-E", "-c", probe).CombinedOutput(); err != nil {
		detail := "validate-skills needs Python 3.10 or later with PyYAML as python3"
		if msg := strings.TrimSpace(string(out)); msg != "" {
			detail += ": " + lastLine(msg)
		}
		return r.Abort(python, "no-python", detail)
	}

	tmp, err := os.MkdirTemp("", "factory-validator")
	if err != nil {
		return r.Abort("validator", "unwritable", err.Error())
	}
	defer os.RemoveAll(tmp)
	scripts := filepath.Join(tmp, "scripts")
	if err := os.Mkdir(scripts, 0o755); err != nil {
		return r.Abort("validator", "unwritable", err.Error())
	}
	if err := os.WriteFile(filepath.Join(scripts, "__init__.py"), nil, 0o644); err != nil {
		return r.Abort("validator", "unwritable", err.Error())
	}
	if err := os.WriteFile(filepath.Join(scripts, "quick_validate.py"), core.QuickValidate, 0o644); err != nil {
		return r.Abort("validator", "unwritable", err.Error())
	}

	for _, skill := range found {
		args := []string{"-B", "-E", "-c", boot, skill, "--codes", "--forbid-tool-names"}
		if strict {
			args = append(args, "--strict")
		}
		cmd := exec.Command(python, args...)
		cmd.Env = append(os.Environ(), "FACTORY_VALIDATOR_DIR="+tmp)
		var stdout, stderr bytes.Buffer
		cmd.Stdout, cmd.Stderr = &stdout, &stderr
		runErr := cmd.Run()

		errorsBefore := r.Errors()
		sc := bufio.NewScanner(&stdout)
		for sc.Scan() {
			rest, ok := strings.CutPrefix(sc.Text(), prefix)
			if !ok {
				continue
			}
			for _, severity := range []string{"error", "warning"} {
				if finding, ok := strings.CutPrefix(rest, severity+" "); ok {
					r.Finding(severity, finding)
				}
			}
		}
		var exitErr *exec.ExitError
		switch {
		case runErr == nil:
		case errors.As(runErr, &exitErr) && exitErr.ExitCode() == 1 && r.Errors() > errorsBefore:
		default:
			detail := lastLine(strings.TrimSpace(stderr.String()))
			if detail == "" {
				detail = runErr.Error()
			}
			r.Error(filepath.ToSlash(skill), 0, "validator-failed", detail)
		}
	}
	return r.Finish(len(found))
}

// find returns dir when it holds a SKILL.md, and otherwise each folder
// below dir that holds one. It does not look below a skill, because a
// second SKILL.md there is a finding of the validator.
func find(dir string) ([]string, error) {
	var found []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}
		if d.Name() == ".git" || d.Name() == "node_modules" {
			return filepath.SkipDir
		}
		if _, err := os.Stat(filepath.Join(path, "SKILL.md")); err == nil {
			found = append(found, path)
			return filepath.SkipDir
		}
		return nil
	})
	return found, err
}

func lastLine(s string) string {
	lines := strings.Split(s, "\n")
	return strings.TrimSpace(lines[len(lines)-1])
}
