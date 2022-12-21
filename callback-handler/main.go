// https://blakewilliams.me/posts/handling-macos-url-schemes-with-go

package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#include "callback-handler.h"
*/
import "C"

import (
	"log"
	"net"
	"net/url"
)

var urlListener chan string = make(chan string)

func main() {
	go C.RunApp()
	urlStr := <-urlListener

	url, err := url.Parse(urlStr)
	if err != nil {
		log.Fatal(err)
	}
	sockAddr := url.Path

	c, err := net.Dial("unix", sockAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	_, err = c.Write([]byte(urlStr))
	if err != nil {
		log.Fatal("write error:", err)
	}
}

//export HandleURL
func HandleURL(u *C.char) {
	urlListener <- C.GoString(u)
}
