package main

import (
	"fmt"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func RunCommand(command string, envvars []string) {
	log.WithFields(log.Fields{"command": command}).Info("Running command")
	cmd := exec.Command(command)
	cmd.Env = append(os.Environ(), envvars...)
	out, err := cmd.Output()
	if err != nil {
		log.Error(err)
	} else {
		fmt.Printf("%s\n", out)
	}
}
