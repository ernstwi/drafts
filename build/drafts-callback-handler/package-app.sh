#!/bin/bash

# Setup
mkdir -p "dist/Drafts CLI Helper.app/Contents/MacOS"
cp "build/drafts-callback-handler/Info.plist" "dist/Drafts CLI Helper.app/Contents/"

cp $BIN "dist/Drafts CLI Helper.app/Contents/MacOS/"
