FROM golang:1.11-alpine AS builder
WORKDIR /go/src/github.com/gmr/env-aws-params
COPY . .
RUN apk add git make\
    && go get -u github.com/golang/dep/cmd/dep \
    && make all

FROM alpine:3.8
COPY --from=builder /go/src/github.com/gmr/env-aws-params/env-aws-params /
RUN apk add --no-cache ca-certificates

ENTRYPOINT [ "/env-aws-params" ]
