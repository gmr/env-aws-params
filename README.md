# env-aws-params
``env-aws-params`` is a tool that injects AWS EC2 Systems Manager (SSM) [Parameter Store](https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-paramstore.html) 
Key / Value pairs as [Environment Variables](https://en.wikipedia.org/wiki/Environment_variable) when executing an
application. It is intended to be used as a Docker [Entrypoint](https://docs.docker.com/engine/reference/builder/#entrypoint),
but can really be used to launch applications outside of Docker as well.

The primary goal is to provide a way of injecting environment variables for
[12 Factor](https://12factor.net) applications that have their configuration defined
in the SSM Parameter store. It was directly inspired by [envconsul](https://github.com/hashicorp/envconsul).

## Example Usage
Create parameters in Parameter Store:
```bash
aws ssm put-parameter --name /service-prefix/ENV_VAR1 --value example
aws ssm put-parameter --name /service-prefix/ENV_VAR2 --value test-value
```

Then use ``env-aws-params`` to have bash display the env vars it was called with:
```bash
env-aws-params --prefix /service-prefix --command /bin/bash -- -c set
```

## Building
This project uses [dep](http://github.com/golang/dep). To build the project:

```bash
dep ensure
go build
```