#!/bin/bash

# Setup
mkdir -p "Drafts Callback Handler.app/Contents/MacOS"
cp "build/drafts-callback-handler/Info.plist" "Drafts Callback Handler.app/Contents/"

cp $BIN "Drafts Callback Handler.app/Contents/MacOS/"
zip -r "drafts-callback-handler-app_$TARGET.zip" "Drafts Callback Handler.app/"
rm -rf "Drafts Callback Handler.app"
