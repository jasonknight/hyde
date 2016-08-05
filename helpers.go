package main

import "fmt"
import "os"
import "strings"

func version() string {
	return "v1.0"
}

func fileExists(p string) bool {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return false
	}
	return true
}
func banner() {
	fmt.Println("*******************************************************")
	fmt.Printf("*                       Hyde %s                     *\n", version())
	fmt.Println("*******************************************************")
}
func FileType(p string) string {
	parts := strings.Split(p, ".")
	ftype := parts[len(parts)-1]
	return ftype
}
func IsCompilable(p string) bool {
	ftypes := []string{"md", "htm", "html", "css", "js", "txt", "csv", "json", "xml"}
	ftype := FileType(p)
	for _, t := range ftypes {
		if ftype == t {
			return true
		}
	}
	return false
}
func IsMarkdown(p string) bool {
	ftype := FileType(p)
	if ftype == "md" || ftype == "markdown" {
		return true
	}
	return false
}
func IsCSS(p string) bool {
	ftype := FileType(p)
	if ftype == "css" {
		return true
	}
	return false
}
func IsHTML(p string) bool {
	ftype := FileType(p)
	if ftype == "htm" || ftype == "html" {
		return true
	}
	return false
}
func IsJS(p string) bool {
	ftype := FileType(p)
	if ftype == "js" {
		return true
	}
	return false
}
