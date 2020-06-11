package main

import (
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// GHToc GitHub TOC
type GHToc []string

// GHDoc GitHub document
type GHDoc struct {
	Path     string
	Depth    int
	Indent   int
	html     string
}

// NewGHDoc create GHDoc
func NewGHDoc(Path string, Depth int, Indent int) *GHDoc {
	return &GHDoc{Path, Depth, Indent, ""}
}

// Convert2HTML downloads remote file
func (doc *GHDoc) Convert2HTML() error {

	if _, err := os.Stat(doc.Path); os.IsNotExist(err) {
		return err
	}

	file, err := os.Open(doc.Path)
	if err != nil {
		return err
	}
	doc.html = Md2Html(file)
	file.Close()	
	return nil
}

// GrabToc gets TOC from html
func (doc *GHDoc) GrabToc() *GHToc {

	re := `(?si)<h(?P<num>[1-6])>\s*` +
//		`<a\s*id="user-content-[^"]*"\s*class="anchor"\s*` +
//		`href="(?P<href>[^"]*)"[^>]*>\s*` +
//		`.*?</a>` +
		`(?P<name>.*?)</h`
	r := regexp.MustCompile(re)
	listIndentation := generateListIndentation(doc.Indent)

	toc := GHToc{}
	minHeaderNum := 6
	var groups []map[string]string
	for idx, match := range r.FindAllStringSubmatch(doc.html, -1) {
		if idx == -1 {
			continue
		}
		group := make(map[string]string)
		// fill map for groups
		for i, name := range r.SubexpNames() {
			if i == 0 || name == "" {
				continue
			}
			group[name] = removeStuff(match[i])
		}
		// update minimum header number
		n, _ := strconv.Atoi(group["num"])
		if n < minHeaderNum {
			minHeaderNum = n
		}
		groups = append(groups, group)
	}

	var tmpSection string
	for _, group := range groups {
		// format result
		n, _ := strconv.Atoi(group["num"])
		if doc.Depth > 0 && n > doc.Depth {
			continue
		}

		tmpSection = removeStuff(group["name"])

		link := strings.Replace(doc.Path,"\\","/",-1)
		link += anchorize(tmpSection)
		
		tocItem := strings.Repeat(listIndentation(), n-minHeaderNum) + "* " +
			"[" + tmpSection + "]" +
			"(" + link + ")"
		toc = append(toc, tocItem)
	}

	return &toc
}

// GetToc return GHToc for a document
func (doc *GHDoc) GetToc() *GHToc {
	if err := doc.Convert2HTML(); err != nil {
		log.Fatal(err)
		return nil
	}
	return doc.GrabToc()
}
