# env-aws-params

[![Build Status](https://travis-ci.org/gmr/env-aws-params.svg?branch=master)](https://travis-ci.org/gmr/env-aws-params)

``env-aws-params`` is a tool that injects AWS EC2 Systems Manager (SSM)
[Parameter Store](https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-paramstore.html)
Key / Value pairs as [Environment Variables](https://en.wikipedia.org/wiki/Environment_variable)
when executing an application. It is intended to be used as a Docker
[Entrypoint](https://docs.docker.com/engine/reference/builder/#entrypoint),
but can really be used to launch applications outside of Docker as well.

The primary goal is to provide a way of injecting environment variables for
[12 Factor](https://12factor.net) applications that have their configuration defined
in the SSM Parameter store. It was directly inspired by
[envconsul](https://github.com/hashicorp/envconsul).

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

## Building

This project uses [go modules](https://go.dev/blog/using-go-modules). To build the project:

```bash
go mod download
go mod verify
go build
```

Building an environment is also provided as a docker image based on Alpine Linux. See the Dockerfile for more information.

```bash
docker build -t env-aws-params; # Build the image
docker run --rm -it -v $HOME/.aws/:/root/.aws/ env-aws-params [your options]
```
