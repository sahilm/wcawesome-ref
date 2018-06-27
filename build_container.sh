#! /usr/bin/env bash

GOOS=linux GOARCH=amd64 go build main.go

docker build -t smuthoo/wcawesome-ref .
