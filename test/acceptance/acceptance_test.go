// Package acceptance runs the acceptance tests of the factory: one
// testscript file per automated scenario, in one folder per increment.
package acceptance

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

var (
	// repo is the root of the factory repo.
	repo string
	// binDir holds the factory command, built from repo.
	binDir string
	// goEnv holds the Go settings that a go command in a script needs,
	// because HOME points to an empty folder there.
	goEnv []string
	// coverDir collects the coverage data that the factory command writes
	// in every script.
	coverDir string
)

// TestMain builds the factory command and installs the stand-in for gh and
// the command stdin-probe, which every script finds on PATH.
func TestMain(m *testing.M) {
	testscript.Main(runner{m}, map[string]func(){
		"gh":          func() { os.Exit(ghStandIn()) },
		"stdin-probe": func() { os.Exit(stdinProbe()) },
	})
}

type runner struct{ m *testing.M }

func (r runner) Run() int { return run(r.m) }

func run(m *testing.M) int {
	var err error
	if repo, err = filepath.Abs(filepath.Join("..", "..")); err != nil {
		fmt.Fprintln(os.Stderr, "acceptance:", err)
		return 2
	}
	if binDir, err = os.MkdirTemp("", "factory-bin"); err != nil {
		fmt.Fprintln(os.Stderr, "acceptance:", err)
		return 2
	}
	defer os.RemoveAll(binDir)
	if coverDir, err = os.MkdirTemp("", "factory-cover"); err != nil {
		fmt.Fprintln(os.Stderr, "acceptance:", err)
		return 2
	}
	defer os.RemoveAll(coverDir)

	// A failed build does not stop the run: each test then fails on its
	// own, because the factory command is missing. The command is built
	// with coverage of the packages of the module.
	build := exec.Command("go", "build", "-cover", "-coverpkg=./...", "-o", binDir+string(os.PathSeparator), "./cmd/factory")
	build.Dir = repo
	if out, err := build.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "acceptance: build factory: %v\n%s", err, out)
	}

	names := []string{"GOCACHE", "GOMODCACHE", "GOPATH", "GOPROXY", "GOFLAGS", "GOTOOLCHAIN"}
	out, err := exec.Command("go", append([]string{"env"}, names...)...).Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, "acceptance: go env:", err)
		return 2
	}
	for i, value := range strings.Split(strings.TrimSuffix(string(out), "\n"), "\n") {
		goEnv = append(goEnv, names[i], value)
	}
	code := m.Run()
	reportCoverage()
	return code
}

// reportCoverage prints the share of statements that the scripts reached,
// and each function that no script reached. A problem here does not change
// the result of the tests.
func reportCoverage() {
	profile := filepath.Join(binDir, "coverage.txt")
	textfmt := exec.Command("go", "tool", "covdata", "textfmt", "-i="+coverDir, "-o="+profile)
	textfmt.Dir = repo
	if out, err := textfmt.CombinedOutput(); err != nil {
		fmt.Printf("acceptance: no coverage: %v %s\n", err, strings.TrimSpace(string(out)))
		return
	}
	module, err := exec.Command("go", "list", "-C", repo, "-m").Output()
	if err != nil {
		fmt.Println("acceptance: no coverage:", err)
		return
	}
	prefix := strings.TrimSpace(string(module)) + "/"
	funcs := exec.Command("go", "tool", "cover", "-func="+profile)
	funcs.Dir = repo
	out, err := funcs.Output()
	if err != nil {
		fmt.Println("acceptance: no coverage:", err)
		return
	}
	var total string
	var uncovered []string
	for _, line := range strings.Split(string(out), "\n") {
		f := strings.Fields(line)
		if len(f) != 3 {
			continue
		}
		if f[0] == "total:" {
			total = strings.TrimSuffix(f[2], "%")
			continue
		}
		if f[2] == "0.0%" {
			place := strings.TrimSuffix(strings.TrimPrefix(f[0], prefix), ":")
			uncovered = append(uncovered, fmt.Sprintf("acceptance: uncovered %s %s", place, f[1]))
		}
	}
	fmt.Printf("acceptance: coverage %s%% of statements\n", total)
	for _, line := range uncovered {
		fmt.Println(line)
	}
}

// TestAcceptance runs every testscript file in the increment folders.
func TestAcceptance(t *testing.T) {
	dirs, err := filepath.Glob("[0-9][0-9]-*")
	if err != nil {
		t.Fatal(err)
	}
	if len(dirs) == 0 {
		t.Fatal("no increment folders")
	}
	for _, dir := range dirs {
		t.Run(dir, func(t *testing.T) {
			testscript.Run(t, testscript.Params{
				Dir:   dir,
				Setup: setup,
				Cmds: map[string]func(*testscript.TestScript, bool, []string){
					"exitcode": exitcode,
					"jsonsub":  jsonsub,
					"jsoneq":   jsoneq,
				},
			})
		})
	}
}

// setup gives each file its own HOME, no system git config, the factory
// command and the stand-in for gh on PATH, REPO, FACTORY_BIN (the folder
// of the factory command alone) and GH_STAND_IN (the folder of the
// stand-in's state), and GOCOVERDIR for the coverage data of the factory
// command. testscript already gives each file a fresh $WORK and
// TMPDIR.
func setup(env *testscript.Env) error {
	home := filepath.Join(env.WorkDir, ".home")
	if err := os.Mkdir(home, 0o755); err != nil {
		return err
	}
	env.Setenv("HOME", home)
	env.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	env.Setenv("REPO", repo)
	env.Setenv("FACTORY_BIN", binDir)
	env.Setenv("GH_STAND_IN", filepath.Join(env.WorkDir, "gh"))
	env.Setenv("GOCOVERDIR", coverDir)
	env.Setenv("PATH", binDir+string(os.PathListSeparator)+env.Getenv("PATH"))
	for i := 0; i < len(goEnv); i += 2 {
		env.Setenv(goEnv[i], goEnv[i+1])
	}
	return nil
}

// exitcode N command [args...] runs command and checks that it exits with
// N. The output stays available to stdout and stderr.
func exitcode(ts *testscript.TestScript, neg bool, args []string) {
	if neg {
		ts.Fatalf("unsupported: ! exitcode")
	}
	if len(args) < 2 {
		ts.Fatalf("usage: exitcode N command [args...]")
	}
	want, err := strconv.Atoi(args[0])
	if err != nil {
		ts.Fatalf("exitcode: bad exit code %q", args[0])
	}
	got := 0
	if err := ts.Exec(args[1], args[2:]...); err != nil {
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) {
			ts.Fatalf("exitcode: %v", err)
		}
		got = exitErr.ExitCode()
	}
	if got != want {
		ts.Fatalf("exit code %d, want %d", got, want)
	}
}
