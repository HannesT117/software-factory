// Command factory checks skills, design notes, decision records, acceptance
// tests and the paths of pull requests, runs the agent hooks, and applies
// the admin setup of a repo.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	claudecode "github.com/HannesT117/software-factory/adapters/claude-code"
	"github.com/HannesT117/software-factory/internal/adr"
	"github.com/HannesT117/software-factory/internal/crossreview"
	"github.com/HannesT117/software-factory/internal/design"
	"github.com/HannesT117/software-factory/internal/gitx"
	"github.com/HannesT117/software-factory/internal/prcheck"
	"github.com/HannesT117/software-factory/internal/report"
	"github.com/HannesT117/software-factory/internal/setup"
	"github.com/HannesT117/software-factory/internal/skills"
	"github.com/HannesT117/software-factory/internal/testcheck"
)

const usage = `factory: usage: factory <command> [args]
factory: validate-skills [--strict] [dir ...]
factory: check-design <increment-dir>
factory: check-adr [dir]
factory: check-tests [--base <ref>]
factory: check-pr --base <ref>
factory: scenarios
factory: cross-review --base <ref>
factory: setup --agent <login> [--repo <owner>/<name>] [--dry-run]
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
		fs := flag.NewFlagSet(cmd, flag.ContinueOnError)
		fs.SetOutput(io.Discard)
		strict := fs.Bool("strict", false, "")
		if err := fs.Parse(args); err != nil {
			return r.Abort("args", "usage", "factory validate-skills [--strict] [dir ...]")
		}
		dirs := fs.Args()
		if len(dirs) == 0 {
			dirs = []string{"."}
		}
		return skills.Validate(r, dirs, *strict)

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

	case "scenarios":
		if len(args) != 0 {
			return r.Abort("args", "usage", "factory scenarios")
		}
		return scenarios(r)

	case "cross-review":
		fs := flag.NewFlagSet(cmd, flag.ContinueOnError)
		fs.SetOutput(io.Discard)
		base := fs.String("base", "", "")
		if err := fs.Parse(args); err != nil || fs.NArg() != 0 || *base == "" {
			return r.Abort("args", "usage", "factory cross-review --base <ref>")
		}
		root, err := gitx.Toplevel(".")
		if err != nil {
			return r.Abort(".", "not-a-repo", err.Error())
		}
		return crossreview.Run(r, root, *base, os.Getenv(crossreview.KeyVar))

	case "setup":
		const setupUsage = "factory setup --agent <login> [--repo <owner>/<name>] [--dry-run]"
		fs := flag.NewFlagSet(cmd, flag.ContinueOnError)
		fs.SetOutput(io.Discard)
		agent := fs.String("agent", "", "")
		repo := fs.String("repo", "", "")
		dryRun := fs.Bool("dry-run", false, "")
		if err := fs.Parse(args); err != nil || fs.NArg() != 0 || *agent == "" {
			return r.Abort("args", "usage", setupUsage)
		}
		if !setup.ValidLogin(*agent) {
			return r.Abort("args", "usage", "--agent takes a GitHub login: "+setupUsage)
		}
		return setup.Run(r, setup.Options{Agent: *agent, Repo: *repo, DryRun: *dryRun, Stdin: stdin})

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
	file := filepath.ToSlash(filepath.Join(dir, "design.md"))
	data, err := os.ReadFile(file)
	if err != nil {
		return r.Abort(file, "unreadable", err.Error())
	}

	// earlier holds the scenarios of the notes beside dir with a lower
	// number, for the replaced scenarios.
	earlier := map[string][]design.Scenario{}
	siblings, _ := filepath.Glob(filepath.Join(filepath.Dir(filepath.Clean(dir)), "*", "design.md"))
	for _, f := range siblings {
		n, ok := design.FolderNumber(filepath.Dir(f))
		if !ok || n >= nn {
			continue
		}
		text, err := os.ReadFile(f)
		if err != nil {
			return r.Abort(filepath.ToSlash(f), "unreadable", err.Error())
		}
		list, _ := design.Scenarios(string(text))
		earlier[n] = append(earlier[n], list...)
	}
	return r.Finish(design.Check(r, file, nn, string(data), earlier))
}

// scenarios prints each scenario in force of the notes under
// docs/increments.
func scenarios(r *report.Report) int {
	root, err := gitx.Toplevel(".")
	if err != nil {
		root = "."
	}
	notes, err := design.ReadNotes(root, "")
	if err != nil {
		return r.Abort(root, "unreadable", err.Error())
	}
	list := design.InForce(notes)
	for _, s := range list {
		r.Line(s.ID + " " + s.Kind + " " + s.Title)
	}
	return r.Finish(len(list))
}
