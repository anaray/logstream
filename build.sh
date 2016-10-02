#!/usr/bin/env bash
rm -f ../../bin/service_*
TARGET=$1
go get github.com/mitchellh/gox
if [ -z $target ]; then
  TARGET="linux/amd64"
fi
gox -osarch=$TARGET -output "../../bin/logstream_{{.OS}}_{{.Arch}}" logstream/service
