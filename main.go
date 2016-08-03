package main

import "flag"
import "fmt"
import "io/ioutil"
import "os"
import "errors"
import "strings"

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
            CompileDirectory(s,strings.Join(np,"/"))
            continue 
        }
        CompileFile(s,strings.Join(np,"/"))
    }


    return nil
}

func CompileFile(s Settings, p string) error {
    fmt.Printf("CompileFile [%s]\n",p)
    return nil
}