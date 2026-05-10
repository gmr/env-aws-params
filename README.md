# env-aws-params

[![CI](https://github.com/gmr/env-aws-params/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/gmr/env-aws-params/actions/workflows/ci.yml)

``env-aws-params`` is a tool that injects AWS EC2 Systems Manager (SSM)
[Parameter Store](https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-paramstore.html)
Key / Value pairs as [Environment Variables](https://en.wikipedia.org/wiki/Environment_variable)
when executing an application. It is intended to be used as a Docker
[Entrypoint](https://docs.docker.com/engine/reference/builder/#entrypoint),
but can really be used to launch applications outside of Docker as well.

If no ``--prefix`` is set, the command is executed with the existing environment
unchanged — ``env-aws-params`` will not contact SSM. This makes it safe to wire
in as a Docker entrypoint that's only active when ``PARAMS_PREFIX`` is set.

The primary goal is to provide a way of injecting environment variables for
[12 Factor](https://12factor.net) applications that have their configuration defined
in the SSM Parameter store. It was directly inspired by
[envconsul](https://github.com/hashicorp/envconsul).

## Installation

### Pre-built binaries

Static binaries for `linux/amd64`, `linux/arm64`, `darwin/amd64`, and `darwin/arm64`
are attached to each [GitHub Release](https://github.com/gmr/env-aws-params/releases).

```bash
VERSION=v1.1.0   # pick from the releases page
OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')

curl -fL "https://github.com/gmr/env-aws-params/releases/download/${VERSION}/env-aws-params_${OS}-${ARCH}" \
  -o /usr/local/bin/env-aws-params
chmod +x /usr/local/bin/env-aws-params
```

### Docker

Multi-arch container images (`linux/amd64`, `linux/arm64`) are published to
the GitHub Container Registry:

```bash
docker pull ghcr.io/gmr/env-aws-params:latest
```

Available tags:

- `latest` — most recent commit on `main`
- `main` — alias for `latest`
- `<version>`, `<major>.<minor>`, `<major>` — tagged releases (e.g. `1.1.0`, `1.1`, `1`)

The image's entrypoint is `env-aws-params`, so flags and the wrapped command go directly to `docker run`:

```bash
docker run --rm \
  -v "$HOME/.aws:/root/.aws:ro" \
  -e AWS_PROFILE \
  ghcr.io/gmr/env-aws-params:latest \
  --prefix /service-prefix /bin/sh -c env
```

Or use it as a base image to wrap your own app:

```dockerfile
FROM ghcr.io/gmr/env-aws-params:latest
COPY my-app /usr/local/bin/my-app
CMD ["my-app"]
```

## Example Usage

Create parameters in Parameter Store:
```bash
aws ssm put-parameter --name /service-prefix/ENV_VAR1 --value example
aws ssm put-parameter --name /service-prefix/ENV_VAR2 --value test-value
```

Then use ``env-aws-params`` to have bash display the env vars it was called with:
```bash
env-aws-params --prefix /service-prefix /bin/bash -c set
```

If you want to include common and service specific values, ``--prefix`` can be specified
multiple times:
```bash
env-aws-params --prefix /common /bin/bash -c set
```

To get a plaintext output of your environment variables to use with other utilities, we can use `printenv`:
```bash
env-aws-params --pristine --silent --prefix /service-prefix /usr/bin/printenv > ~/some-file.sh
```
Which will write your environment variables in plain text, for example:
```bash
# ~/some-file.sh Contents:
ENV_VAR1=example
ENV_VAR2=test-value
```

## CLI Options

```
NAME:
   env-aws-params - Application entry-point that injects SSM Parameter Store values as Environment Variables

USAGE:
   env-aws-params [global options] -p prefix command [command arguments]

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --aws-region value        The AWS region to use for the Parameter Store API [$AWS_REGION]
   --prefix value, -p value  Key prefix that is used to retrieve the environment variables - supports multiple use
   --pristine                Only use values retrieved from Parameter Store, do not inherit the existing environment variables
   --sanitize                Replace invalid characters in keys to underscores
   --strip                   Strip invalid characters in keys
   --upcase                  Force keys to uppercase
   --debug                   Log additional debugging information [$PARAMS_DEBUG]
   --silent                  Silence all logs [$PARAMS_SILENT]
   --help, -h                show help
   --version, -v             print the version
```

## Building from source

This project uses [Go modules](https://go.dev/blog/using-go-modules) and requires Go 1.26+.

```bash
go mod download
go build
```

Or build the Docker image locally:

```bash
docker build -t env-aws-params .
```
