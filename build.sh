#!/bin/sh
#apk add --update go
mkdir -p /go/src/github.com/Dataman-Cloud/omega-es
mkdir -p $HOME/.omega/
export GOPATH=/go
export GO15VENDOREXPERIMENT=1

cp -r /src/* /go/src/github.com/Dataman-Cloud/omega-es
rm /etc/localtime && cd /go/src/github.com/Dataman-Cloud/omega-es && mv localtime /etc
cd /go/src/github.com/Dataman-Cloud/omega-es && mv start.sh /bin/ && mv omega-es.yaml $HOME/.omega/ && cd src && go build && mv src /bin/omega-es
apk del go
rm -rf /go
rm -rf /src
rm -rf /var/cache/apk/*
