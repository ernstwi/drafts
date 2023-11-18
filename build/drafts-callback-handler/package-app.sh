#!/bin/bash

# Run from root of repo

# Create .app structure
mkdir -p "Drafts Callback Handler.app/Contents/MacOS"
cp "dist/drafts-callback-handler_darwin_arm64/drafts-callback-handler" "Drafts Callback Handler.app/Contents/MacOS/"
cp "build/drafts-callback-handler/Info.plist" "Drafts Callback Handler.app/Contents/"

# Zip it
zip -r "drafts-callback-handler.zip" "Drafts Callback Handler.app/"
rm -r "Drafts Callback Handler.app"
