// +build !windows,!plan9

package term

import (
	"fmt"
	"os"

	"github.com/elves/elvish/pkg/diag"
	"github.com/elves/elvish/pkg/sys"
	"golang.org/x/sys/unix"
)

func setup(in, out *os.File) (func() error, error) {
	// On Unix, use input file for changing termios. All fds pointing to the
	// same terminal are equivalent.

	fd := int(in.Fd())
	term, err := sys.TermiosForFd(fd)
	if err != nil {
		return nil, fmt.Errorf("can't get terminal attribute: %s", err)
	}

	savedTermios := term.Copy()

	term.SetICanon(false)
	term.SetIExten(false)
	term.SetEcho(false)
	term.SetVMin(1)
	term.SetVTime(0)

	// Enforcing crnl translation on readline. Assuming user won't set
	// inlcr or -onlcr, otherwise we have to hardcode all of them here.
	term.SetICRNL(true)

	err = term.ApplyToFd(fd)
	if err != nil {
		return nil, fmt.Errorf("can't set up terminal attribute: %s", err)
	}

	var errSetupVT error
	err = setupVT(out)
	if err != nil {
		errSetupVT = fmt.Errorf("can't setup VT: %s", err)
	}

	restore := func() error {
		return diag.Errors(savedTermios.ApplyToFd(fd), restoreVT(out))
	}

	return restore, errSetupVT
}

func setupGlobal() func() {
	return func() {}
}

func sanitize(in, out *os.File) {
	// Some programs use non-blocking IO but do not correctly clear the
	// non-blocking flags after exiting, so we always clear the flag. See #822
	// for an example.
	unix.SetNonblock(int(in.Fd()), false)
	unix.SetNonblock(int(out.Fd()), false)
}
