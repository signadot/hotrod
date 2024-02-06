FROM golang:1.21-alpine

ARG TARGETPLATFORM

COPY dist/$TARGETPLATFORM/bin/hotrod /app/hotrod

ENTRYPOINT ["/app/hotrod"]

