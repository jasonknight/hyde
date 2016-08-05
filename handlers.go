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
					return v.url
				}
			}
			return ""
		}
		fmt.Println("Registered link_to")
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
