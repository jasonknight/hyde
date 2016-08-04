package main

import "fmt"
import "errors"
import "strings"
import "io/ioutil"
import "os"
import "bytes"
import "text/template"
import "bufio"

func CompileDirectory(s Settings, p string) error {
    fmt.Printf("CompileDirectory [%s]\n",p)
    if ( ! fileExists(p) ) {
        return errors.New(fmt.Sprintf("%s does not exist",p))
    }
    err := discoverLayout(&s,p)
    if (err != nil) {
        fmt.Println(err)
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
    dfile := ConvertPath(s,dpath,"html")
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
    tmpl, err := template.New(name).Funcs(s.fmap).Parse(text)

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

func CompileGoPartial(s Settings, p string) (string,error) {
     file_contents,err := ioutil.ReadFile(p)
    if ( err != nil ) {

        return "",err
    }
    tmpl, err := template.New(p).Funcs(s.fmap).Parse(string(file_contents[:]))

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