FROM golang:1.13-alpine AS builder
MAINTAINER Totoval <totoval@tobyan.com> (https://totoval.com)

LABEL "com.github.actions.name"="Go Release Binary"
LABEL "com.github.actions.description"="Automate publishing Totoval build artifacts for GitHub releases"
LABEL "com.github.actions.icon"="cpu"
LABEL "com.github.actions.color"="blue"

LABEL "name"="Automate publishing Totoval build artifacts for GitHub releases through GitHub Actions"
LABEL "version"="1.0.0"
LABEL "repository"="http://github.com/totoval/go-release.action"
LABEL "homepage"="https://totoval.com"

LABEL "maintainer"="Totoval <totoval@tobyan.com> (https://totoval.com)"

ENV GO111MODULE on

ADD . /src

WORKDIR /src

RUN go mod download
RUN go build -o ./builds/artisan ./artisan.go

# ===================================================================================
FROM scratch
MAINTAINER Totoval <totoval@tobyan.com> (https://totoval.com)
LABEL "maintainer"="Totoval <totoval@tobyan.com> (https://totoval.com)"

COPY --from=builder /src/builds/artisan /sms/sms-forwarder
COPY --from=builder /src/.env.example.json /sms/.env.json

WORKDIR /sms

ENTRYPOINT /sms/sms-forwarder sms:read /dev/ttyUSB0