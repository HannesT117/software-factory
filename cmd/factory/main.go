// Command factory checks skills, design notes, decision records, acceptance
// tests and the paths of pull requests, and runs the agent hooks.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	claudecode "github.com/HannesT117/software-factory/adapters/claude-code"
	"github.com/HannesT117/software-factory/internal/adr"
	"github.com/HannesT117/software-factory/internal/design"
	"github.com/HannesT117/software-factory/internal/gitx"
	"github.com/HannesT117/software-factory/internal/prcheck"
	"github.com/HannesT117/software-factory/internal/report"
	"github.com/HannesT117/software-factory/internal/skills"
	"github.com/HannesT117/software-factory/internal/testcheck"
)

const usage = `factory: usage: factory <command> [args]
factory: validate-skills [dir ...]
factory: check-design <increment-dir>
factory: check-adr [dir]
factory: check-tests [--base <ref>]
factory: check-pr --base <ref>
factory: hook pre-tool-use
factory: hook stop
`

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout))
}

func run(args []string, stdin io.Reader, stdout io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(stdout, usage)
		return report.Usage
	}
	cmd, args := args[0], args[1:]
	r := report.New(cmd, stdout)
	switch cmd {
	case "validate-skills":
		if len(args) == 0 {
			args = []string{"."}
		}
		return skills.Validate(r, args)

	case "check-design":
		if len(args) != 1 {
			return r.Abort("args", "usage", "factory check-design <increment-dir>")
		}
		return checkDesign(r, args[0])

	case "check-adr":
		dir := "docs/decisions"
		if len(args) > 1 {
			return r.Abort("args", "usage", "factory check-adr [dir]")
		}
		if len(args) == 1 {
			dir = args[0]
		}
		return adr.Check(r, dir)

	case "check-tests":
		fs := flag.NewFlagSet(cmd, flag.ContinueOnError)
		fs.SetOutput(io.Discard)
		base := fs.String("base", "", "")
		if err := fs.Parse(args); err != nil || fs.NArg() != 0 {
			return r.Abort("args", "usage", "factory check-tests [--base <ref>]")
		}
		root, err := gitx.Toplevel(".")
		if err != nil {
			if *base != "" {
				return r.Abort(".", "not-a-repo", err.Error())
			}
			root = "."
		}
		return testcheck.Check(r, root, *base)

	case "check-pr":
		fs := flag.NewFlagSet(cmd, flag.ContinueOnError)
		fs.SetOutput(io.Discard)
		base := fs.String("base", "", "")
		if err := fs.Parse(args); err != nil || fs.NArg() != 0 || *base == "" {
			return r.Abort("args", "usage", "factory check-pr --base <ref>")
		}
		return prcheck.Run(r, ".", *base)

	case "hook":
		if len(args) != 1 {
			return r.Abort("args", "usage", "factory hook pre-tool-use|stop")
		}
		switch args[0] {
		case "pre-tool-use":
			return claudecode.PreToolUse(stdin, stdout)
		case "stop":
			return claudecode.Stop(stdin, stdout)
		}
		return r.Abort("args", "usage", "factory hook pre-tool-use|stop")
	}
	fmt.Fprint(stdout, usage)
	return report.Usage
}

func checkDesign(r *report.Report, dir string) int {
	nn, ok := design.FolderNumber(dir)
	if !ok {
		return r.Abort(dir, "usage", "the folder name is not <nn>-<slug>")
	}
	file := dir + "/design.md"
	data, err := os.ReadFile(file)
	if err != nil {
		return r.Abort(file, "unreadable", err.Error())
	}
	return r.Finish(design.Check(r, file, nn, string(data)))
}
