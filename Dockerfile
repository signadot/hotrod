FROM golang:1.27.1-alpine

ARG TARGETPLATFORM

COPY dist/$TARGETPLATFORM/bin/hotrod /app/hotrod

ENTRYPOINT ["/app/hotrod"]

