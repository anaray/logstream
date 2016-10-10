#!/usr/bin/env bash
#rm -f ../../bin/service_*
TARGET=$1
glide update
go get github.com/mitchellh/gox
if [ -z $target ]; then
  gox -osarch=$TARGET -output "$GOPATH/bin/logstream_{{.OS}}_{{.Arch}}" github.com/anaray/logstream/main
fi
