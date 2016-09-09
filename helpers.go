package main

import "fmt"
import "os"
import "strings"
import "io"
//import "io/ioutil"
import "bufio"

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
func MoveFile(s Settings, p string) error {
	dpath, err := MakeDestinationPath(s, p)
	dfile := ConvertPath(s, dpath, FileType(p))
	in, err := os.Open(p)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dfile)
	if err != nil {
		return err
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	err = out.Sync()
	return err
}
func MakeDestinationPath(s Settings, p string) (string, error) {
	fmt.Printf("MakeDestinationPath received: %s\n",p)
	var ps []string
	var tps []string
	var fpath string
	ps = strings.Split(p, "/")
	ps[0] = s.outdir
	tps = ps[0 : len(ps)-1]
	fpath = strings.Join(tps, "/")
	fmt.Printf("About to make: %s\n",fpath)
	err := os.MkdirAll(fpath, 0777)
	if err != nil {
		return "", err
	}
	return strings.Join(ps, "/"), nil
}

func DestinationPath(s Settings, p string) string {
	var ps []string
	ps = strings.Split(p, "/")
	ps[0] = s.outdir
	if IsMarkdown(p) {
		return ConvertPath(s, strings.Join(ps, "/"), "html")
	} else {
		return ConvertPath(s, strings.Join(ps, "/"), FileType(p))
	}

}

func DestinationURL(s Settings, p string) string {
	var ps []string
	if IsMarkdown(p) {
		ps = strings.Split(ConvertPath(s, p, "html"), "/")
	} else {
		ps = strings.Split(ConvertPath(s, p, FileType(p)), "/")
	}

	ps[0] = s.url
	return strings.Join(ps, "/")
}

func DestinationURLPath(s Settings, p string) string {
	var ps []string
	if IsMarkdown(p) {
		ps = strings.Split(ConvertPath(s, p, "html"), "/")
	} else {
		ps = strings.Split(ConvertPath(s, p, FileType(p)), "/")
	}

	return "/" + strings.Join(ps[1:],"/")
}

func ConvertPath(s Settings, p string, t string) string {
	parts := strings.Split(p, ".")
	parts[len(parts)-1] = t
	return strings.Join(parts, ".")
}
func FilePutContents(p string,txt string) error {
	f, err := os.Create(p)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(f)
	_, err = w.WriteString(txt)
	w.Flush()
	return nil
}
func GenerateFile(p, t string) {
	if t == "md" {
		err := FilePutContents(p,DefaultMDTemplate())
		if err != nil {
			panic(err)
		}
		return
	} 
	fmt.Printf("Didn't know what to do with %s\n",t)
}