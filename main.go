package main

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	app := cli.NewApp()
	app.Name = "env-aws-params"
	app.Usage = "Application entry-point that injects SSM Parameter Store values as Environment Variables"
	app.UsageText = "env-aws-params [global options] -p prefix command [command arguments]"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "aws-region",
			Usage:  "The AWS region to use for the Parameter Store API",
			EnvVar: "AWS_REGION",
		},
		cli.StringSliceFlag{
			Name:  "prefix, p",
			Usage: "Key prefix that is used to retrieve the environment variables - supports multiple use",
		},
		cli.BoolFlag{
			Name:  "pristine",
			Usage: "Only use values retrieved from Parameter Store, do not inherit the existing environment variables",
		},
		cli.BoolFlag{
			Name:  "sanitize",
			Usage: "Replace invalid characters in keys to underscores",
		},
		cli.BoolFlag{
			Name:  "upcase",
			Usage: "Force keys to uppercase",
		},
	}

	app.Action = func(c *cli.Context) error {
		var envvars []string

		if len(c.GlobalStringSlice("prefix")) == 0 {
			log.Fatal("prefix is required")
			os.Exit(1)
		}

		fmt.Print(c.GlobalString("aws-region"))

		client, err := NewSSMClient(c.GlobalString("aws-region"))
		if err != nil {
			log.Fatal(err)
			os.Exit(2)
		}

		for _, path := range c.GlobalStringSlice("prefix") {
			params, err := client.GetParametersByPath(path)
			if err != nil {
				log.Error(err)
				os.Exit(3)
			}
			vars := BuildEnvVars(params, c.GlobalBool("upcase"))
			envvars = append(envvars, vars...)
		}

		if c.GlobalBool("pristine") == false {
			envvars = append(os.Environ(), envvars...)
		}

		RunCommand(c.Args()[0], c.Args()[1:], envvars)
		return nil
	}

	app.Run(os.Args)
}

func BuildEnvVars(parameters map[string]string, upcase bool) []string {
	var vars []string

	for k, v := range parameters {
		if upcase == true {
			k = strings.ToUpper(k)
		}
		vars = append(vars, fmt.Sprintf("%s=%s", k, v))
	}
	return vars
}
