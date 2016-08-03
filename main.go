package main

import "flag"
import "fmt"
import "io/ioutil"
import "os"
import "errors"
import "strings"
import "text/template"
import "bytes"

type Settings struct {
    indir string
    outdir string
}
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
    fmt.Printf("*                       Hyde %s                     *\n",version())
    fmt.Println("*******************************************************")
}

func main() {
    var (
        indir = flag.String("in","./_src","The location of all the files")
        outdir = flag.String("out","./_dest", "Where to put the output html")

    )
    
    flag.Parse()
    s := Settings{indir: *indir, outdir: *outdir}
    banner()
    fmt.Printf("In %s and out %s\n", s.indir, s.outdir)
    fmt.Println(DefaultLayout())
    err := CompileDirectory(s, s.indir)
    if ( err != nil ) {
        panic(err)
    }
    fmt.Println("Done.")
}

func CompileDirectory(s Settings, p string) error {
    fmt.Printf("CompileDirectory [%s]\n",p)
    if ( ! fileExists(p) ) {
        return errors.New(fmt.Sprintf("%s does not exist",p))
    }
    flist,err := ioutil.ReadDir(p)
    if ( err != nil ) {
        return err
    }

    for _,f := range flist {
        np := []string{p,f.Name()}
        if ( f.IsDir() ) { 
            err = CompileDirectory(s,strings.Join(np,"/"))
            if ( err != nil ) {
                return err
            }
            continue 
        }
        err = CompileFile(s,strings.Join(np,"/"))
        if ( err != nil ) {
            return err
        }
    }


    return nil
}

func CompileFile(s Settings, p string) error {
    fmt.Printf("CompileFile [%s]\n",p)
    dpath,err := MakeDestinationPath(s,p)
    if ( err != nil ) {
        return err
    }
    dfile,_ := ConvertPath(s,dpath,"html")
    fmt.Printf("Destination: %s => %s\n",dpath, dfile)

    compiled,err := CompileGoTemplate(s,p)
    fmt.Println(compiled)
    return nil
}

func MakeDestinationPath(s Settings, p string) (string,error) {
    var ps []string
    var tps []string
    var fpath string
    ps = strings.Split(p,"/")
    ps[0] = s.outdir
    tps = ps[0:len(ps) - 1]
    fpath = strings.Join(tps,"/")
    err := os.MkdirAll(fpath,0777)
    if ( err !=  nil ) {
        return "",err
    }
    return strings.Join(ps,"/"),nil
}

func ConvertPath(s Settings, p string, t string) (string,error) {
    parts := strings.Split(p,".")
    parts[len(parts)-1] = t
    return strings.Join(parts,"."),nil
}

func CompileGoTemplate(s Settings, p string) (string, error) {
    file_contents,err := ioutil.ReadFile(p)

    tmpl, err := template.New(p).Parse(string(file_contents[:]))

    if ( err != nil) {
        return "",err
    }
    var can bytes.Buffer
    err = tmpl.Execute(&can,s)
    if (err != nil) {
        return "",err
    }
    return can.String(),nil
}