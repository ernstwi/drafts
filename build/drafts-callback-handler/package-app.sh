#!/bin/bash

# Setup
mkdir -p "Drafts CLI Helper.app/Contents/MacOS"
cp "build/drafts-callback-handler/Info.plist" "Drafts CLI Helper.app/Contents/"

cp $BIN "Drafts CLI Helper.app/Contents/MacOS/"
