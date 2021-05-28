FROM ubuntu
COPY ./hotrod /go/bin/hotrod
ENTRYPOINT ["/go/bin/hotrod"]
