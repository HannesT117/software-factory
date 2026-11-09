package acceptance

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
)

// stdinProbeData is what stdin-probe puts on the standard input of the
// command. It ends in a newline, so that a command that reads one line
// gets a value and does not stop at an empty input.
const stdinProbeData = "stdin-probe: this input must stay unread\n"

// stdinProbe runs "stdin-probe <command> [args...]". It gives the command a
// pipe as standard input that holds stdinProbeData, with the write end
// already closed, and passes on its output and exit code. After the command
// exits, it reads what is left in the pipe. If any byte is missing, the
// command read standard input: stdin-probe says so on stderr and exits 99.
// The check needs no timeout: the input is in the pipe before the command
// starts, and a read returns at once, with data or with the end of input.
func stdinProbe() int {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: stdin-probe command [args...]")
		return 2
	}
	r, w, err := os.Pipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "stdin-probe:", err)
		return 2
	}
	if _, err := io.WriteString(w, stdinProbeData); err != nil {
		fmt.Fprintln(os.Stderr, "stdin-probe:", err)
		return 2
	}
	w.Close()
	defer r.Close()

	cmd := exec.Command(os.Args[1], os.Args[2:]...)
	cmd.Stdin = r
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	code := 0
	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) {
			fmt.Fprintln(os.Stderr, "stdin-probe:", err)
			return 2
		}
		code = exitErr.ExitCode()
	}

	left, err := io.ReadAll(r)
	if err != nil {
		fmt.Fprintln(os.Stderr, "stdin-probe:", err)
		return 2
	}
	if read := len(stdinProbeData) - len(left); read > 0 {
		fmt.Fprintf(os.Stderr, "stdin-probe: %s read %d bytes of standard input\n", os.Args[1], read)
		return 99
	}
	return code
}
