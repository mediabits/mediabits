#!/bin/sh
#
# DO NOT USE THIS SCRIPT MANUALLY
# Use go generate

cd "$GOPATH"/src/mediabits

rm -f assets/assets.go
go-bindata \
	-o assets/assets.go \
	-pkg assets \
	-prefix 'assets/static/' \
	-nomemcopy \
	assets/static/templates \
	assets/static/html_templates \
	assets/static/html_static