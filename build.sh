#! /bin/bash

go build -ldflags "-s -w"
npm i  --no-package-lock
npm run build
