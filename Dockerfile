FROM alpine:3.8 as builder

RUN apk add --no-cache go dep git make gcc musl-dev
RUN mkdir -p /opt/go/src/env-aws-params.git \
    && cd /opt/go/src/ \
    && git clone --depth=1 https://github.com/gmr/env-aws-params.git
RUN /bin/sh -c 'cd /opt/go/src/env-aws-params && export GOPATH=/opt/go && make all'


FROM alpine:3.8
COPY --from=builder /opt/go/src/env-aws-params/env-aws-params /
RUN apk add --no-cache ca-certificates

ENTRYPOINT [ "/env-aws-params" ]

