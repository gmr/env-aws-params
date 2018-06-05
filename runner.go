package main

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func RunCommand(command string, args []string, envVars []string) error {

	log.Infof("PID %v running %s %s", os.Getpid(), command,
		strings.Join(args[:], " "))

	procAttr := new(os.ProcAttr)
	procAttr.Env = envVars
	procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}

	// prefix args with the command, as per https://golang.org/pkg/os/#StartProcess
	// The argv slice will become os.Args in the new process, so it normally starts
	// with the program name.
	args = append([]string{command}, args...)
	proc, err := os.StartProcess(command, args, procAttr)
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
		switch sigv {
		case syscall.SIGHUP:
			err = syscall.Kill(-os.Getpid(), syscall.SIGHUP)
		case syscall.SIGINT:
			err = syscall.Kill(-os.Getpid(), syscall.SIGINT)
		case syscall.SIGTERM:
			err = syscall.Kill(-os.Getpid(), syscall.SIGTERM)
		case syscall.SIGQUIT:
			err = syscall.Kill(-os.Getpid(), syscall.SIGQUIT)
		default:
			err = syscall.Kill(-os.Getpid(), syscall.SIGTERM)
		}
		log.WithFields(log.Fields{
			"err":    err,
			"proc":   proc,
			"pid":    -proc.Pid,
			"signal": sigv},
		).Info("Caught signal, sent to child")
	}()
	_, err = proc.Wait()
	log.Debug("Exiting")
	return err
}
