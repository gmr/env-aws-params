package main

import (
	"fmt"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func RunCommand(command string, args []string, envvars []string) {
	log.WithFields(log.Fields{
		"command": command,
		"args":    args},
	).Info("Running command")
	cmd := exec.Command(command, args...)
	cmd.Env = envvars
	out, err := cmd.Output()
	if err != nil {
		log.Error(err)
	} else {
		fmt.Printf("%s\n", out)
	}
}
