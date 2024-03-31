// https://blakewilliams.me/posts/handling-macos-url-schemes-with-go

package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#include "callback-handler.h"
*/
import "C"

import (
	"errors"
	"log"
	"net"
	"net/url"
	"os/exec"
	"strings"
)

var urlListener chan string = make(chan string)

func fatal(err error) {
	if err != nil {
		// TODO: Pass error message to drafts-cli output
		log.Fatal(err)
	}
}

//export HandleURL
func HandleURL(u *C.char) {
	urlListener <- C.GoString(u)
}

func main() {
	go C.RunApp()
	urlStr := <-urlListener

	// Drafts does not properly escape ";" and "+"
	// https://en.wikipedia.org/wiki/URL_encoding#Percent-encoding_reserved_characters
	urlStr = strings.ReplaceAll(urlStr, ";", "%3B")
	urlStr = strings.ReplaceAll(urlStr, "+", "%2B")

	url, err := url.Parse(urlStr)
	fatal(err)
	sockAddr := url.Path

	c, err := net.Dial("unix", sockAddr)
	fatal(err)
	defer c.Close()

	_, err = c.Write([]byte(urlStr))
	fatal(err)

	v := url.Query()
	if !v.Has("app") {
		fatal(errors.New("missing `app` URL parameter"))
	}
	err = exec.Command("open", "-a", v.Get("app")).Run()
	fatal(err)
}
