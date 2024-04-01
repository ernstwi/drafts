# Drafts CLI

Command line interface for [Drafts](https://getdrafts.com). Requires Drafts Pro.

https://github.com/ernstwi/drafts/assets/17042548/3fe84265-f30e-480e-9e8a-25feecde1b04

## Install

1. Install the required [helper action](https://directory.getdrafts.com/a/2Qx) from the Drafts Directory

2. Install the CLI:

```
brew install ernstwi/tap/drafts
```

3. Install the required helper application (see ["Implementation notes"](#implementation-notes) below):

```
brew install --cask --no-quarantine ernstwi/tap/drafts-helper
```

## Usage

```
$ drafts --help
Usage: drafts <command> [<args>]

Options:
  --help, -h             display this help and exit

Commands:
  new                    create new draft
  prepend                prepend to draft
  append                 append to draft
  replace                append to draft
  edit                   edit draft in $EDITOR
  get                    get content of draft
  select                 select active draft using fzf
```

See further: `drafts <command> --help`

## Implementation notes

It is easy to send commands _to_ Drafts. This is done using actions defined in Drafts' [URL scheme](https://docs.getdrafts.com/docs/automation/urlschemes). Using a custom action, we can send arbitrary JavaScript to be executed in Drafts.

The more difficult thing is how to get a response back from Drafts to the CLI. Drafts' URL scheme offers a `x-success` parameter which can be used to call a separate URL with the result of an action. The problem is that the `drafts` CLI process can't catch this call directly.

My solution to this problem was to make a separate macOS app (Drafts CLI Helper) whose only purpose is to catch `x-success` calls from Drafts, and forward them via Unix socket to the appropriate `drafts` CLI process. The resulting system is illustrated below.

```mermaid
graph LR;
    cli[Drafts CLI]-- URL scheme -->drafts[Drafts]
    drafts-- URL scheme -->cbh[Drafts CLI Helper]
    cbh-- Unix socket -->cli
```

## Dev setup

```
ln -s "$PWD/dist/Drafts CLI Helper.app" /Applications/
PATH=$PATH:$PWD/dist/drafts_darwin_all
./build.sh   # To build
./release.sh # To build and publish a new release
```
