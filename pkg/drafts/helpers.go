package drafts

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func open(action string, v url.Values) url.Values {
	ch := make(chan string)
	go server(ch)
	sockAddr := <-ch // Wait for ready signal
	v.Add("x-success", "drafts-callback-handler://"+sockAddr)
	err := exec.Command("open", "-g", draftsURL(action, v)).Run()
	fatal(err)

	return urlValues(<-ch)
}

func urlValues(urlstr string) url.Values {
	u, err := url.Parse(urlstr)
	fatal(err)
	return u.Query()
}

func draftsURL(action string, v url.Values) string {
	return fmt.Sprintf("drafts://x-callback-url/%s?%s", action, strings.ReplaceAll(v.Encode(), "+", "%20"))
}

// Start a server, listen for one message, send it over ch
func server(ch chan string) {
	// Create a temp file to use as socket address
	f, err := os.CreateTemp("", "*.sock")
	fatal(err)
	sockAddr := f.Name()

	// We don't actually need the file, just the filename
	err = f.Close()
	fatal(err)
	os.Remove(f.Name())

	// To delete the socket after communication is done
	defer os.Remove(f.Name())

	l, err := net.Listen("unix", sockAddr)
	fatal(err)
	defer l.Close()

	// Send socket address to open(), which sends it to the callback handler via
	// Drafts. The callback handler then uses the socket address to forward the
	// reply from Drafts to open(). This also signals to open() that the server
	// is ready to accept connections.
	ch <- sockAddr

	c, err := l.Accept()
	fatal(err)

	msg, err := io.ReadAll(c)
	fatal(err)

	ch <- string(msg)
}

func mustJSON(a any) string {
	js, err := json.Marshal(a)
	fatal(err)
	return string(js)
}
