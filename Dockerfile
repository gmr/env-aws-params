ARG GO_VERSION=1.26

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS builder

ARG TARGETOS
ARG TARGETARCH
ARG VERSION=dev

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -trimpath -ldflags "-w -s -X main.VersionString=${VERSION}" \
    -o /out/env-aws-params .

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=builder /out/env-aws-params /usr/local/bin/env-aws-params
ENTRYPOINT ["/usr/local/bin/env-aws-params"]
