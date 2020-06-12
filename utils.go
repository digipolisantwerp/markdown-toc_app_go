package main

import (
	"strings"
	"os"
	"regexp"
	"io/ioutil"
	"net/url"
)

var (
//	cleanReg = regexp.MustCompile(`[^a-zA-Z0-9\-\s\%]+`)
	noPunctuationReg = regexp.MustCompile(`[.,?!'\";:\&]+`)
)

// removeStuff trims spaces, removes new lines and code tag from a string
func removeStuff(s string) string {
	res := strings.Replace(s, "\n", "", -1)
	res = strings.Replace(res, "<code>", "", -1)
	res = strings.Replace(res, "</code>", "", -1)
	res = strings.TrimSpace(res)

	return res
}

// generate func of custom spaces indentation
func generateListIndentation(spaces int) func() string {
	return func() string {
		return strings.Repeat(" ", spaces)
	}
}

func readFile2String(file string) string {
	fileContent, err := ioutil.ReadFile(*tocFile)
	if err != nil {
		return ""
	}
	return string(fileContent)
}

func writeString2File(content string, file string) {
	f, err := os.Create(file)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	f.WriteString(content)
}

func anchorize(link string) string {
	return "#" + strings.ToLower(strings.ReplaceAll(url.QueryEscape(noPunctuationReg.ReplaceAllString(link,"")),"+","-"))
}

func stringInSlice(val string, list []string) bool {
    for _, b := range list {
        if strings.EqualFold(b,val) {
            return true
		}
    }
    return false
}