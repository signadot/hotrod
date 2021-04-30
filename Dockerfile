FROM ubuntu
ARG TARGETARCH=amd64
COPY hotrod-linux-$TARGETARCH /go/bin/hotrod-linux
ENTRYPOINT ["/go/bin/hotrod-linux"]