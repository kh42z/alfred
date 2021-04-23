#!/usr/bin/env bash

set -e
cd "${0%/*}/../.."


echo "> go fmt"
gofmt -s -w .
