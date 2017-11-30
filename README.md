# env-aws-params
Application entry-point that injects [AWS SSM Parameter Store](https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-paramstore.html) 
Key / Value pairs as [Environment Variables](https://en.wikipedia.org/wiki/Environment_variable).

## Example Usage

```bash
env-aws-params --prefix /path/to/kv/pairs --command /bin/bash -- -c set
```

## Building
This project uses [dep](http://github.com/golang/dep). To build the project:

```bash
dep ensure
go build
```