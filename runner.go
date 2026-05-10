package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
)

type CommandFailedError struct {
	ExitCode int
}

func (e *CommandFailedError) Error() string {
	return fmt.Sprintf("Command failed with exit code %d", e.ExitCode)
}

func RunCommand(command string, args []string, envVars []string) error {
	// Resolve command against $PATH so callers can pass bare names like "ls"
	// rather than absolute paths. Absolute / relative paths pass through
	// unchanged after an existence + executable-bit check.
	resolved, err := exec.LookPath(command)
	if err != nil {
		return err
	}

	log.Infof("PID %v running %s %s", os.Getpid(), resolved,
		strings.Join(args, " "))

	procAttr := new(os.ProcAttr)
	procAttr.Env = envVars
	procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}

	// prefix args with the command, as per https://golang.org/pkg/os/#StartProcess
	// The argv slice will become os.Args in the new process, so it normally starts
	// with the program name.
	args = append([]string{resolved}, args...)
	proc, err := os.StartProcess(resolved, args, procAttr)
	if err != nil {
		return err
	}
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		sigv := <-sigc
		var killErr error
		switch sigv {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			killErr = syscall.Kill(-os.Getpid(), sigv.(syscall.Signal))
		default:
			killErr = syscall.Kill(-os.Getpid(), syscall.SIGTERM)
		}
		log.WithFields(log.Fields{
			"err":    killErr,
			"proc":   proc,
			"pid":    -proc.Pid,
			"signal": sigv},
		).Info("Caught signal, sent to child")
	}()
	procState, err := proc.Wait()
	if err != nil {
		return err
	}
	if procState.ExitCode() != 0 {
		return &CommandFailedError{procState.ExitCode()}
	}
	return nil
}
