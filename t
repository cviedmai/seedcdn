#!/bin/bash
go test ./... -test.run="."
rc=$?
if [[ $rc != 0 ]] ; then
  exit $rc
fi
