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
)

func TestMain(m *testing.M) {
	os.Exit(run(m))
}

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

	// A failed build does not stop the run: each test then fails on its
	// own, because the factory command is missing.
	build := exec.Command("go", "build", "-o", binDir+string(os.PathSeparator), "./cmd/factory")
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
	return m.Run()
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
				Cmds:  map[string]func(*testscript.TestScript, bool, []string){"exitcode": exitcode},
			})
		})
	}
}

// setup gives each file its own HOME, no system git config, the factory
// command on PATH, and REPO. testscript already gives each file a fresh
// $WORK and TMPDIR.
func setup(env *testscript.Env) error {
	home := filepath.Join(env.WorkDir, ".home")
	if err := os.Mkdir(home, 0o755); err != nil {
		return err
	}
	env.Setenv("HOME", home)
	env.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	env.Setenv("REPO", repo)
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
