#!/bin/sh

# Remove any old artifacts
rm -rf build/
rm -rf target/

# Cross-compile to all supported platforms
goxc

# Rename windows binaries with an exe extension
find build -iwholename '*windows*' -type f -exec mv {} {}.exe \;

# Generate the update JSON file
"$GOPATH"/src/mediabits-updater/mediabits-updater

# Copy everything to the target
mkdir -p target
cp -r build target/updates
cp docs/index.html target/index.html