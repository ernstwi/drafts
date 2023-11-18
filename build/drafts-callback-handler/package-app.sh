#!/bin/bash

# Run from root of repo

function setup() {
	mkdir -p "Drafts Callback Handler.app/Contents/MacOS"
	cp "build/drafts-callback-handler/Info.plist" "Drafts Callback Handler.app/Contents/"
}

rm -f *.zip

setup
cp "dist/drafts-callback-handler_darwin_amd64_v1/drafts-callback-handler" "Drafts Callback Handler.app/Contents/MacOS/"
zip -r "drafts-callback-handler-intel.zip" "Drafts Callback Handler.app/"
rm -rf "Drafts Callback Handler.app"

setup
cp "dist/drafts-callback-handler_darwin_arm64/drafts-callback-handler" "Drafts Callback Handler.app/Contents/MacOS/"
zip -r "drafts-callback-handler-arm.zip" "Drafts Callback Handler.app/"
rm -rf "Drafts Callback Handler.app"
