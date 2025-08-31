#!/bin/bash

cd ./cmd/server

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../../iot-backend
