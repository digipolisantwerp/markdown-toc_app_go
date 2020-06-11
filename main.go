package main

import (
	"os"
	"strings"
	"regexp"
	"path/filepath"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	version   = "1.0.0"
)

var (
    tocFile     = kingpin.Flag("toc-file", "File that will contain TOC").Default("README.md").String()
	startFolder = kingpin.Flag("start-folder", "Starting folder of MD files").Default(".").String()
	depth       = kingpin.Flag("depth", "How many levels of headings to include. Defaults to 0 (all)").Default("0").Int()
	indent      = kingpin.Flag("indent", "Indent space of generated list").Default("2").Int()
)

// Entry point
func main() {
	kingpin.Version(version)
	kingpin.Parse()

    var paths []string

    err := filepath.Walk(*startFolder, func(path string, info os.FileInfo, err error) error {
		//only .md files
		if strings.EqualFold(filepath.Ext(path),".md") {
			// if md file is not our toc file
			if !strings.EqualFold(path,*tocFile) {
				paths = append(paths, path)
			}
		}
        return nil
    })
    if err != nil {
        panic(err)
    }

	pathsCount := len(paths)

	ch := make(chan *GHToc, pathsCount)

	for _, path := range paths {
		ghdoc := NewGHDoc(path, *depth, *indent)
		ch <- ghdoc.GetToc()
	}

	combinedToc := ""
	for i := 1; i <= pathsCount; i++ {
		tocStr := ""
		toc := <-ch
		if toc != nil {
			for _, tocItem := range *toc {
				tocStr += tocItem
				tocStr += "\n"
			}
//			tocStr += "\n"
		}
		combinedToc += tocStr
	}

	// update TOC in file
	outputToc := "<!-- dgp-toc-start -->\n"
	outputToc += combinedToc
	outputToc += "<!-- dgp-toc-end -->"

	if _, err := os.Stat(*tocFile); err == nil {
		fileContent := readFile2String(*tocFile)
		tocReg := regexp.MustCompile(`(?s)<!-- dgp-toc-start -->.*<!-- dgp-toc-end -->`)
		outputToc = string(tocReg.ReplaceAll([]byte(fileContent), []byte(outputToc)))
	}
	
	writeString2File(outputToc, *tocFile)
}