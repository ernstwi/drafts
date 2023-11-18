#!/bin/bash

# cd to script directory
cd "$(dirname "$0")"

# Create .app structure
mkdir -p "Drafts Callback Handler.app/Contents/MacOS"
cp "../../dist/drafts-callback-handler_darwin_arm64/drafts-callback-handler" "Drafts Callback Handler.app/Contents/MacOS/"
cp "Info.plist" "Drafts Callback Handler.app/Contents/"

# Zip it
zip -r "Drafts Callback Handler.zip" "Drafts Callback Handler.app/"
rm -r "Drafts Callback Handler.app"
