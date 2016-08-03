package main

import "flag"
import "fmt"
import "io/ioutil"
import "os"
import "errors"
import "strings"
import "text/template"
import "bytes"
import "bufio"
import "regexp"

type Settings struct {
    indir string
    outdir string
    layout string
    prepends []string
    file_ids map[string]string
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
    s.file_ids = make(map[string]string)
    banner()
    fmt.Printf("In %s and out %s\n", s.indir, s.outdir)
    s.layout = DefaultLayout()
    err := discoverLayout(&s)
    if (err != nil) {
        fmt.Println(err)
    }
    
    err = discoverFileIds(&s,s.indir)
    if (err != nil) {
        fmt.Println(err)
    }
    //fmt.Println(s.layout)
    fmt.Println("Beginning Compilation")
    err = CompileDirectory(s, s.indir)
    if ( err != nil ) {
        panic(err)
    }
    fmt.Println("Done.")
}
func discoverLayout(s *Settings) error {
    flist,err := ioutil.ReadDir(s.indir)
    if ( err != nil ) {
        return err
    }
    for _,f := range flist {
        np := []string{s.indir,f.Name()}
        if ( f.IsDir() ) { 
            continue 
        }
        if (f.Name() == "_layout.md" ) {
            file_contents,err := ioutil.ReadFile(strings.Join(np,"/"))
            if ( err != nil ) {
                return err
            }
            s.layout = string(file_contents[:])
            return nil
        }
    }
    return errors.New("No Layout Found")
}
func discoverFileIds(s *Settings, p string) error {
    //fmt.Printf("discoverFileIds[%s]\n",p)
    if ( ! fileExists(p) ) {
        return errors.New(fmt.Sprintf("%s does not exist",p))
    }
    flist,err := ioutil.ReadDir(p)
    if ( err != nil ) {
        return err
    }

    for _,f := range flist {
        np := []string{p,f.Name()}
        fname := f.Name()
        if ( f.IsDir()) { 
            err = discoverFileIds(s,strings.Join(np,"/"))
            if ( err != nil ) {
                return err
            }
            continue 
        }
        r, err := regexp.Compile("^([\\w\\d]+)--.+")
        if ( err != nil ) {
            return err
        }
        matches := r.FindStringSubmatch(fname)
        //fmt.Println(matches)
        if ( len(matches) >= 2 ) {
            //fmt.Printf("id: [%s]\n",matches[1])
            s.file_ids[ matches[1] ] = strings.Join(np,"/")
        } 
        
    }
    return nil
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
        fname := f.Name()
        if (fname[0] == '_') {
            continue
        }
        if ( f.IsDir()) { 
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

    if ( err != nil ) {
        return err
    }
    if ( compiled == "" ) {
        return nil
    }
    f, err := os.Create(dfile)
    if ( err != nil) {
        return err
    }
    w := bufio.NewWriter(f)
    _, err = w.WriteString(compiled)
    w.Flush()
    if ( err != nil ) {
        return err
    }
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
    if ( err != nil ) {
        return "",err
    }
    return CompileGoString(s,p,string(file_contents[:]))
    
}
func CompileGoString(s Settings,name string, text string) (string,error) {
    // First we parse the string for special directives
    var flines []string
    oflines := strings.Split(text,"\n")
    for _,line := range oflines {
        if ( len(line) > 2 && line[0] == '#' && line[1] == '!' ) {
            l2exe := line[3:len(line)]
            l2exe = fmt.Sprintf("{{%s}}",l2exe)
            fmt.Println("Prepending ", l2exe)
            s.prepends = append(s.prepends,l2exe)
            continue
        }
        flines = append(flines,line)
    }
    text = strings.Join(flines,"\n")
    text = strings.Replace(s.layout,"{{.Content}}",text,1)
    text = strings.Join(s.prepends,"\n") + "\n" + text
    //fmt.Println("Text to compile is: ", text)
    tmpl, err := template.New(name).Parse(text)

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
