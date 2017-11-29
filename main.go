package main

import (
	"os"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)


func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}


func main() {
	var command string
	var prefix string

	app := cli.NewApp()
	app.Name = "env-aws-params"
	app.Usage = "Application entry-point that injects SSM Parameter Store values as Environment Variables"

	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "prefix, p",
			Value: "",
			Usage: "Key prefix that is used to retrieve the environment variables ",
			Destination: &prefix,
		},
		cli.StringFlag{
			Name: "command, c",
			Value: "",
			Usage: "Command",
			Destination: &command,
		},
	}

	app.Action = func(c *cli.Context) error {
		var vars []string

		ssm, err := NewSSMClient()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		parameters, err := ssm.getParameters(prefix)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		for k, v := range parameters {
			vars = append(vars, fmt.Sprintf("%s=%s", k, v))
		}

		RunCommand(command, vars)
		return nil
	}

	app.Run(os.Args)
}
