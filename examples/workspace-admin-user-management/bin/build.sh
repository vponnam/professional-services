#!/bin/bash

case "$1" in 
"linux")
  echo "Building for linux"
  env GOOS=linux GOARCH=amd64 go build -o rotate-linux ../main.go
  ;;
"darwin")
  echo "Building for osx"
  env GOOS=darwin GOARCH=amd64 go build -o rotate-darwin ../main.go
  ;;
  *)
  echo "Usage: build.sh linux/darwin"
  ;;
esac
