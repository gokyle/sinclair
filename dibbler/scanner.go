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
	Date   int64    `json:"title"`
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

	var body = []byte{}
	for scanner.Scan() {
		body = append(body, scanner.Bytes()...)
	}
	body = bytes.TrimSpace(body)

	node.Body = string(blackfriday.MarkdownCommon(body))
	node.Path = path
	return node, nil
}

func returnNodes() {
	var container struct {
		Success bool    `json:"success"`
		Message string  `json:"message"`
		Nodes   []*Node `json:"nodes"`
	}

	if flag.NArg() == 0 {
		container.Message = "no files specified"
	}

	for _, path := range flag.Args() {
		node, err := loadFile(path)
		if err != nil {
			container.Message = err.Error()
			break
		}
		container.Nodes = append(container.Nodes, node)
	}
	if container.Message == "" {
		container.Success = true
	}

	out, err := json.Marshal(container)
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

	var container struct {
		Success bool      `json:"success"`
		Message string    `json:"message"`
		Results []*Result `json:"results"`
	}

	if flag.NArg() == 0 {
		container.Message = "no files specified"
	}

	for _, path := range flag.Args() {
		fi, err := os.Stat(path)
		if err != nil {
			container.Message = err.Error()
			break
		}
		result := &Result{
			Path:  path,
			Mtime: fi.ModTime().Unix(),
		}
		container.Results = append(container.Results, result)
	}
	if container.Message == "" {
		container.Success = true
	}

	out, err := json.Marshal(container)
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
