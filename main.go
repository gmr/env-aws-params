package main

import (
	"context"
	"errors"
	"io"
	"maps"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
	"golang.org/x/sync/errgroup"
)

var VersionString string

func main() {
	cmd := &cli.Command{
		Name:      "env-aws-params",
		Usage:     "Application entry-point that injects SSM Parameter Store values as Environment Variables",
		UsageText: "env-aws-params [global options] -p prefix command [command arguments]",
		Version:   VersionString,
		Flags:     cliFlags(),
		Action:    action,
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func action(ctx context.Context, cmd *cli.Command) error {
	if cmd.Bool("debug") {
		log.SetLevel(log.DebugLevel)
	}
	if cmd.Bool("silent") {
		log.SetOutput(io.Discard)
	}

	code, err := validateArgs(cmd.NArg(), cmd.Bool("sanitize"), cmd.Bool("strip"))
	if code > 0 {
		return cli.Exit(errorPrefix(err), code)
	}

	var envVars []string
	if len(cmd.StringSlice("prefix")) > 0 {
		params, err := getParameters(ctx, cmd)
		if err != nil {
			return cli.Exit(errorPrefix(err), -1)
		}

		envVars = BuildEnvVars(
			params,
			cmd.Bool("sanitize"),
			cmd.Bool("strip"),
			cmd.Bool("upcase"))

		for _, v := range envVars {
			key, _, _ := strings.Cut(v, "=")
			log.Debugf("Setting %s", key)
		}
	} else {
		log.Warn("No prefix set; executing command without retrieving SSM parameters")
	}

	// SSM-derived values come first so they win over inherited environ:
	// glibc/musl/Apple libc all return the first match from getenv().
	if !cmd.Bool("pristine") {
		envVars = append(envVars, os.Environ()...)
	}

	args := cmd.Args()
	if err := RunCommand(args.First(), args.Tail(), envVars); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return cli.Exit(errorPrefix(err), exitErr.ExitCode())
		}
		return cli.Exit(errorPrefix(err), 128)
	}
	return nil
}

func cliFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "aws-region",
			Usage:   "The AWS region to use for the Parameter Store API",
			Sources: cli.EnvVars("AWS_REGION"),
		},
		&cli.StringFlag{
			Name:    "profile",
			Usage:   "Optional AWS profile to use for the Parameter Store API",
			Sources: cli.EnvVars("AWS_PROFILE"),
		},
		&cli.StringSliceFlag{
			Name:    "prefix",
			Aliases: []string{"p"},
			Usage:   "Key prefix that is used to retrieve the environment variables - supports multiple use",
			Sources: cli.EnvVars("PARAMS_PREFIX"),
		},
		&cli.BoolFlag{
			Name:    "pristine",
			Usage:   "Only use values retrieved from Parameter Store, do not inherit the existing environment variables",
			Sources: cli.EnvVars("PARAMS_PRISTINE"),
		},
		&cli.BoolFlag{
			Name:    "sanitize",
			Usage:   "Replace invalid characters in keys to underscores",
			Sources: cli.EnvVars("PARAMS_SANITIZE"),
		},
		&cli.BoolFlag{
			Name:    "strip",
			Usage:   "Strip invalid characters in keys",
			Sources: cli.EnvVars("PARAMS_STRIP"),
		},
		&cli.BoolFlag{
			Name:    "upcase",
			Usage:   "Force keys to uppercase",
			Sources: cli.EnvVars("PARAMS_UPCASE"),
		},
		&cli.BoolFlag{
			Name:    "debug",
			Usage:   "Log additional debugging information",
			Sources: cli.EnvVars("PARAMS_DEBUG"),
		},
		&cli.BoolFlag{
			Name:    "silent",
			Usage:   "Silence all logs",
			Sources: cli.EnvVars("PARAMS_SILENT"),
		},
	}
}

func errorPrefix(err error) string {
	return "ERROR: " + err.Error()
}

func getParameters(ctx context.Context, cmd *cli.Command) (map[string]string, error) {
	client, err := NewSSMClient(ctx, cmd.String("aws-region"), cmd.String("profile"))
	if err != nil {
		return nil, err
	}

	prefixes := cmd.StringSlice("prefix")
	results := make([]map[string]string, len(prefixes))

	g, gctx := errgroup.WithContext(ctx)
	for i, path := range prefixes {
		g.Go(func() error {
			m, err := client.GetParametersByPath(gctx, path)
			if err != nil {
				return err
			}
			results[i] = m
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}

	// Later prefixes override earlier ones when keys collide.
	values := make(map[string]string)
	for _, m := range results {
		maps.Copy(values, m)
	}
	return values, nil
}

func validateArgs(nargs int, sanitize, strip bool) (int, error) {
	if nargs == 0 {
		return 1, errors.New("command not specified")
	}

	if sanitize && strip {
		return 2, errors.New("--sanitize and --strip are mutually exclusive behaviors")
	}

	return 0, nil
}
