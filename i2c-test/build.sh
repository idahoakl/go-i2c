#!/usr/bin/env bash
set -x
GOARM=7 GOARCH=arm GOOS=linux go build -o arm7_i2c-test
