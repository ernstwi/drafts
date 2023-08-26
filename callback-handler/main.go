// https://blakewilliams.me/posts/handling-macos-url-schemes-with-go

package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#include "callback-handler.h"
*/
import "C"

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var urlListener chan string = make(chan string)

func fatal(err error) {
	if err != nil {
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

	termApp := "Terminal"
	{
		conf := config()
		if conf.Terminal != "" {
			termApp = conf.Terminal
		}
	}

	err = exec.Command("open", "-a", termApp).Run()
	fatal(err)
}

type Config struct {
	Terminal string `json:"terminal"`
}

func config() Config {
	home, _ := os.UserHomeDir()
	file, _ := os.Open(filepath.Join(home, ".config", "drafts-cli", "config.json"))
	defer file.Close()
	decoder := json.NewDecoder(file)
	c := Config{}
	err := decoder.Decode(&c)
	if err != nil {
		fmt.Println("error:", err)
	}
	return c
}
