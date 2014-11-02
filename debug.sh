#!/bin/sh
#
# Sets up an asset debugging environment

cd "$GOPATH"/src/mediabits

rm -f assets/assets.go
go-bindata \
	-debug \
	-o assets/assets.go \
	-pkg assets \
	-prefix 'assets/static/' \
	-nomemcopy \
	assets/static/templates \
	assets/static/html_templates \
	assets/static/html_static