#!/bin/bash

# Create .app structure
mkdir -p "DraftsCallbackHandler.app/Contents/MacOS"
cp "../../dist/drafts-callback-handler_darwin_arm64/drafts-callback-handler" "DraftsCallbackHandler.app/Contents/MacOS/"
cp "Info.plist" "DraftsCallbackHandler.app/Contents/"

# Zip it
zip -r "DraftsCallbackHandler.zip" "DraftsCallbackHandler.app/"
