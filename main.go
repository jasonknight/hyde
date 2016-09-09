package main

import "flag"
import "fmt"
import "io/ioutil"
import "errors"
import "strings"

import "regexp"

type FileEntry struct {
	name string
	src  string
	dest string
	url  string
	id   string
	path string
}

type Settings struct {
	indir    string
	outdir   string
	url      string
	layout   string
	prepends []string
	file_ids map[string]FileEntry
	fmap     map[string]interface{}
	HydeMsg  string
}

func init() {
	settingsFilters = make(map[string]SettingsFilter)
}

func main() {
	var (
		indir       = flag.String("in", "./_src", "The location of all the files")
		outdir      = flag.String("out", "./_dest", "Where to put the output html")
		action      = flag.String("action", "compile", "compile|routes")
		url         = flag.String("url", "http://localhost", "the url of your site")
		show_layout = flag.Bool("layout", false, "echo the default template")
		gen_page	= flag.Bool("g",false,"Generate a file")
		gen_path 	= flag.String("p","./_src/un-named.md","The file to generate")
		gen_type	= flag.String("t","md","the type to generate")
	)

	flag.Parse()
	if *show_layout == true {
		fmt.Println(DefaultLayout())
		return
	}
	
	if *gen_page == true {
		GenerateFile(*gen_path,*gen_type);
		return
	}

	s := Settings{indir: *indir, outdir: *outdir, url: *url}
	if s.outdir[len(s.outdir)-1] == '/' || s.indir[len(s.indir)-1] == '/' {
		fmt.Printf("\n\nNo trailing slashes in paths! \n\n")
		return
	}
	s.file_ids = make(map[string]FileEntry)
	s.fmap = make(map[string]interface{})
	banner()
	fmt.Printf("In %s and out %s\n", s.indir, s.outdir)
	s.layout = DefaultLayout()
	s.HydeMsg = HydeMsg()
	err := discoverFileIds(&s, s.indir)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(s.layout)
	registerLinkTo()
	registerPartial()
	for _, fn := range settingsFilters {
		fn(&s)
	}
	if *action == "compile" {
		fmt.Println("Beginning Compilation")
		err = CompileDirectory(s, s.indir)
		if err != nil {
			panic(err)
		}
	}
	if *action == "routes" {
		fmt.Println("Displaying Routes")
		printRoutes(s)
	}

	fmt.Println("Done.")
}
func printRoutes(s Settings) {
	for k, v := range s.file_ids {
		fmt.Printf("link_to(\"%s\") => %s %v\n", k, v.url, v)
	}
}
func discoverLayout(s *Settings, d string) error {
	flist, err := ioutil.ReadDir(d)
	if err != nil {
		return err
	}
	for _, f := range flist {
		np := []string{d, f.Name()}
		if f.IsDir() {
			continue
		}
		if f.Name() == "_layout.html" {
			file_contents, err := ioutil.ReadFile(strings.Join(np, "/"))
			if err != nil {
				return err
			}
			fmt.Println("Discovered layout for ", d)
			s.layout = string(file_contents[:])

			return nil
		}
	}
	return errors.New("No Layout Found")
}
func discoverFileIds(s *Settings, p string) error {
	//fmt.Printf("discoverFileIds[%s]\n",p)
	if !fileExists(p) {
		return errors.New(fmt.Sprintf("%s does not exist", p))
	}
	flist, err := ioutil.ReadDir(p)
	if err != nil {
		return err
	}

	for _, f := range flist {
		np := []string{p, f.Name()}
		fname := f.Name()
		if f.IsDir() {
			err = discoverFileIds(s, strings.Join(np, "/"))
			if err != nil {
				return err
			}
			continue
		}
		r, err := regexp.Compile("^([\\w\\d]+)--.+")
		if err != nil {
			return err
		}
		matches := r.FindStringSubmatch(fname)
		//fmt.Println(matches)
		if len(matches) >= 2 {
			//fmt.Printf("id: [%s]\n",matches[1])
			s.file_ids[matches[1]] = FileEntry{
				name: fname,
				src:  strings.Join(np, "/"),
				dest: DestinationPath(*s, strings.Join(np, "/")),
				id:   matches[1],
				url:  DestinationURL(*s, strings.Join(np, "/")),
				path: DestinationURLPath(*s,strings.Join(np, "/")),
			}
		} else {
			src := strings.Join(np, "/")
			src_split := strings.Split(src, "/")
			fid := strings.Join(src_split[1:], "/")
			s.file_ids[fid] = FileEntry{
				name: fname,
				src:  src,
				dest: DestinationPath(*s, strings.Join(np, "/")),
				id:   fid,
				url:  DestinationURL(*s, strings.Join(np, "/")),
				path: DestinationURLPath(*s,strings.Join(np, "/")),
			}

		}

	}
	return nil
}


