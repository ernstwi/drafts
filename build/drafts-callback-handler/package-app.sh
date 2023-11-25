#!/bin/bash

# Setup
mkdir -p "Drafts Callback Handler.app/Contents/MacOS"
cp "build/drafts-callback-handler/Info.plist" "Drafts Callback Handler.app/Contents/"

cp $BIN "Drafts Callback Handler.app/Contents/MacOS/"
