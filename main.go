// Clu fetches documentation from MDN given a key word and outputs plain text
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"
)

type Dir struct {
	name   string
	url    string
	regexp string
}

func (d Dir) String() string {
	return fmt.Sprintf("%s (%s): %s", d.name, d.url, d.regexp)
}

func main() {
	// TODO: Parse flags

	// Get search query string from either command line or stdin
	var q string
	if len(os.Args) < 2 {
		var err error
		r := bufio.NewReader(os.Stdin)
		q, err = r.ReadString('\n')
		if err != nil {
			log.Fatal("no input given")
		}
	} else {
		q = os.Args[1]
	}
	if q == "" || q == "\n" {
		log.Fatal("no input given")
	}

	// Set up libraries and regexpes to look for
	var dir []Dir
	dir = append(dir, Dir{
		name:   "CSS",
		url:    "https://developer.mozilla.org/en-US/docs/Web/CSS/",
		regexp: `(?s:<article.*</article>)`})
	dir = append(dir, Dir{
		name:   "HTML",
		url:    "https://developer.mozilla.org/en-US/docs/Web/HTML/Element/",
		regexp: `(?s:<article.*</article>)`})

	// Go fetch

	c1 := fetch(dir[0], q)
	c2 := fetch(dir[1], q)

	for done := false; !done; {
		select {
		case msg := <-c1:
			if msg != "" {
				fmt.Printf("%s\n", msg)
				done = true
			}
		case msg := <-c2:
			if msg != "" {
				fmt.Printf("%s\n", msg)
				done = true
			}
		case <-time.After(3 * time.Second):
			done = true
		}
	}
}

func fetch(d Dir, q string) <-chan string {
	c := make(chan string)

	go func() {
		uri := d.url + q
		answer := uri + "\n\n"

		resp, err := http.Get(uri)
		if err != nil {
			log.Print("response error:", err)
		} else {
			defer resp.Body.Close()
			r := bufio.NewReader(resp.Body)
			str, err := r.ReadBytes('\000')
			if err != nil && err != io.EOF {
				log.Print("body read error:", err)
			}
			rex := regexp.MustCompile(d.regexp)
			matched := rex.Find(str)

			if len(matched) > 0 {
				// TODO: Parse HTML into plain text
				c <- answer + string(matched)
			}
		}
		c <- ""
	}()

	return c
}
