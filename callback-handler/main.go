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
)

var urlListener chan string = make(chan string)

//export HandleURL
func HandleURL(u *C.char) {
	urlListener <- C.GoString(u)
}

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

	termApp := "Terminal"
	{
		conf := config()
		if conf.Terminal != "" {
			termApp = conf.Terminal
		}
	}

	err = exec.Command("open", "-a", termApp).Run()
	if err != nil {
		log.Fatal(err)
	}
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
