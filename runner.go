package main

import (
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func RunCommand(command string, args []string, envVars []string) error {
	cmd := exec.Command(command, args...)

	// A nil Env makes os/exec hand the child the parent's environment, which
	// would silently defeat --pristine when there is nothing to inject.
	if envVars == nil {
		envVars = []string{}
	}
	cmd.Env = envVars
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}

	log.Infof("PID %v running %s %s", cmd.Process.Pid, cmd.Path, strings.Join(args, " "))

	sigc := make(chan os.Signal, 32)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		for sigv := range sigc {
			killErr := syscall.Kill(-os.Getpid(), sigv.(syscall.Signal))
			log.WithFields(log.Fields{
				"err":    killErr,
				"pid":    -cmd.Process.Pid,
				"signal": sigv,
			}).Info("Caught signal, sent to child")
		}
	}()

	return cmd.Wait()
}
