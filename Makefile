all: env-aws-params

deps:
	@ echo "Running dep ensure"
	@ /usr/bin/env dep ensure

env-aws-params: deps
	@ echo "Running go build"
	GOOS=linux GOARCH=amd64 go build -ldflags="-w -s -X main.VersionString=v${TRAVIS_TAG}"
