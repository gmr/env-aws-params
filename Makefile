all: env-aws-params

deps:
	/go/bin/dep ensure

env-aws-params: deps
	GOOS=linux GOARCH=amd64 go build -ldflags="-w -s -X main.VersionString=v${TRAVIS_TAG} -a -installsuffix cgo
