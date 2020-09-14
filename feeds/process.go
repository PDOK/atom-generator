package feeds

import (
	"net/http"

	"github.com/imdario/mergo"
)

// ProcessFeed func
func ProcessFeed(f Feed) Feed {
	d := GetDefaultFeedProperties()
	mergo.Merge(&f, d)

	links := f.Link

	var self, describedby, search, up Link
	if f.Self != nil {
		self = Self(*f.Self)
		links = append(links, self)
	}
	if f.Describedby != nil {
		describedby = DescribedBy(*f.Describedby)
		links = append(links, describedby)
	}
	if f.Search != nil {
		search = Search(*f.Search)
		links = append(links, search)
	}
	if f.Up != nil {
		up = Up(*f.Up)
		links = append(links, up)
	}

	f.Link = links

	for i, l := range f.Link {
		f.Link[i] = l.SetHrefLang(*f.Lang)
	}

	for _, e := range f.Entry {
		for j, l := range e.Link {
			if l.Data != nil {
				n := l
				res, err := http.Head(*l.Data)
				if err != nil {
					panic(err)
				}
				if len(l.Length) == 0 {
					n.Length = res.Header.Get("Content-Length")
				}
				if len(l.Type) == 0 {
					n.Type = res.Header.Get("Content-Type")
				}

				n.Data = nil
				e.Link[j] = n
			}
			e.Link[j] = l.SetHrefLang(*f.Lang)
		}
	}

	// reset predefined
	f.Self = nil
	f.Describedby = nil
	f.Search = nil
	f.Up = nil

	return f
}
