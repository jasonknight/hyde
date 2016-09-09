package main

import "fmt"

type SettingsFilter func(s *Settings)

var settingsFilters map[string]SettingsFilter

func RegisterSettingsFilter(name string, f SettingsFilter) {
	settingsFilters[name] = f
}
func registerLinkTo() {
	RegisterSettingsFilter("link_to", func(s *Settings) {
		s.fmap["link_to"] = func(name string) string {
			for k, v := range s.file_ids {
				if k == name {
					return v.path
				}
			}
			return ""
		}
		fmt.Println("Registered link_to")
	})
}
func registerLink() {
	RegisterSettingsFilter("link", func(s *Settings) {
		s.fmap["link"] = func(txt string, name string, id string, cls string) string {
			for k, v := range s.file_ids {
				if k == name {
					return fmt.Sprintf("<a href=\"%s\" id=\"%s\" class=\"%s\">%s</a>",v.path,id,cls,txt)
				}
			}
			return ""
		}
		fmt.Println("Registered link")
	})
}
func registerPartial() {
	RegisterSettingsFilter("partial", func(s *Settings) {
		s.fmap["partial"] = func(name string) string {
			for _, v := range s.file_ids {
				//fmt.Printf("name: %s k: %s id: %s\n",name,k,v.id)
				if name == v.id {
					txt, err := CompileGoPartial(*s, v.src)
					if err == nil {
						return txt
					}
					panic(err)
				} else {
					fmt.Printf("%s != %s\n", name, v.id)
				}
			}

			return ""
		}
		fmt.Println("Registered partial")
	})
}

func registerHTMLHelpers() {
	simple_tags := []string{"div","span","section","h1","h2","h3","h4","h5"}
	for _,tag := range simple_tags {
		RegisterSettingsFilter(tag, func(s *Settings) {
			s.fmap[tag] = func(id string, cls string, innertext string) string {
				return fmt.Sprintf("<%s id=\"%s\" class=\"%s\">%s</%s>",tag,id,cls,innertext,tag)
			}
			fmt.Println("Registered HTML Helper: ",tag)
		})
	}
	RegisterSettingsFilter("a", func(s *Settings) {
		s.fmap["a"] = func(id string, cls string, href string, innertext string) string {
			return fmt.Sprintf("<a id=\"%s\" class=\"%s\" href=\"%s\">%s</a>",id,cls, href, innertext)
		}
		fmt.Println("Registered HTML Helper: a")
	})
	
}