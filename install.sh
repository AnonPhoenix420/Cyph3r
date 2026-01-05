#!/bin/sh
set -e


go mod tidy
go build -o cyph3r ./cmd/cyph3r
sudo install -m 755 cyph3r /usr/local/bin/cyph3r


echo "Installed cyph3r to /usr/local/bin"
