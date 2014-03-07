package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/russross/blackfriday"
	"os"
	"strings"
	"time"
)

var dateFormat = "2006-01-02 15:04"

// Type Node contains the information required to build a page.
type Node struct {
	Title  string   `json:"title"`
	Date   int64    `json:"date"`
	Static bool     `json:"static"`
	Tags   []string `json:"tags"`
	Mtime  int64    `json:"mtime"`
	Path   string   `json:"path"`
	Slug   string   `json:"slug"`
	Body   string   `json:"body"`
}

// loadPost loads the file, extracting metadata and processing the body.
func loadFile(path string) (node *Node, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(file)
	lines := 0
	node = &Node{
		Mtime: fi.ModTime().Unix(),
	}

	for scanner.Scan() {
		line := scanner.Text()
		if lines == 0 {
			if line != "---" {
				return nil, fmt.Errorf("no metadata supplied")
			}
			lines++
			continue
		} else if line == "---" {
			break
		}

		metadata := strings.SplitN(line, ":", 2)
		for i, md := range metadata {
			metadata[i] = strings.TrimSpace(md)
		}

		switch metadata[0] {
		case "title":
			if len(metadata) == 1 {
				return nil, fmt.Errorf("no title specified")
			}
			node.Title = metadata[1]
		case "date":
			if len(metadata) == 1 {
				return nil, fmt.Errorf("invalid time stamp")
			}
			var t time.Time
			t, err = time.Parse(dateFormat, metadata[1])
			if err != nil {
				return nil, err
			}
			node.Date = t.Unix()
		case "static":
			if len(metadata) == 1 {
				return nil, fmt.Errorf("invalid static value")
			}
			if strings.ToLower(metadata[1]) == "true" {
				node.Static = true
			}
		case "tags":
			if len(metadata) == 1 {
				return nil, fmt.Errorf("no tags specified")
			}
			node.Tags = strings.Split(metadata[1], ",")
			for i, t := range node.Tags {
				node.Tags[i] = strings.TrimSpace(t)
			}
		case "slug":
			if len(metadata) == 1 {
				return nil, fmt.Errorf("no slug specified")
			}
			node.Slug = strings.ToLower(metadata[1])
		}
	}

	if lines == 0 {
		return nil, fmt.Errorf("empty file")
	}

	var body = []byte{}
	for scanner.Scan() {
		body = append(body, scanner.Bytes()...)
		body = append(body, 0xa)
	}

	body = bytes.TrimSpace(body)
	if len(body) == 0 {
		return nil, fmt.Errorf("empty file")
	}

	node.Body = strings.TrimSpace(string(blackfriday.MarkdownCommon(body)))
	node.Path = path
	return node, nil
}

func returnNodes() {
	type response struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Path    string `json:"path"`
		Node    *Node  `json:"node"`
	}

	var containers []response

	if flag.NArg() == 0 {
		var resp = response{
			Message: "no files specified",
		}
		containers = append(containers, resp)
	}

	for _, path := range flag.Args() {
		var resp response
		resp.Path = path
		node, err := loadFile(path)
		if err != nil {
			resp.Message = err.Error()
		} else {
			resp.Success = true
			resp.Node = node
		}
		containers = append(containers, resp)
	}

	out, err := json.Marshal(containers)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(string(out))
	}
}

func returnModTime() {
	type Result struct {
		Path  string `json:"path"`
		Mtime int64  `json:"mtime"`
	}

	type response struct {
		Success bool    `json:"success"`
		Message string  `json:"message"`
		Path    string  `json:"path"`
		Result  *Result `json:"result"`
	}

	var containers []response

	if flag.NArg() == 0 {
		resp := response{
			Message: "no files specified",
		}
		containers = append(containers, resp)
	}

	for _, path := range flag.Args() {
		var resp response
		resp.Path = path
		fi, err := os.Stat(path)
		if err != nil {
			resp.Message = err.Error()
		} else {
			result := &Result{
				Path:  path,
				Mtime: fi.ModTime().Unix(),
			}
			resp.Result = result
			resp.Success = true
		}
		containers = append(containers, resp)
	}

	out, err := json.Marshal(containers)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(string(out))
	}
}

func main() {
	mtime := flag.Bool("mod", false, "only retrieve modification times")
	flag.Parse()

	if !*mtime {
		returnNodes()
	} else {
		returnModTime()
	}
}
