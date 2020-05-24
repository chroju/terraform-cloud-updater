FROM golang:1.14.3-alpine

RUN ["/bin/sh", "-c", "apk add --update --no-cache bash ca-certificates curl git jq openssh"]
COPY . /workdir

WORKDIR /workdir
ENTRYPOINT [ "./entrypoint.sh" ]
