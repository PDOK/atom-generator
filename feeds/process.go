package feeds

import (
	"net/http"

	"github.com/imdario/mergo"
)

// ProcessFeed func
func ProcessFeeds(fs Feeds) []Feed {
	var processedFeeds []Feed
	for _, f := range fs.Feeds {
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

		f.recentUpdated(fs)

		for _, entry := range f.Entry {
			for linkIndex, link := range entry.Link {
				if link.Data != nil {
					res, err := http.Head(*link.Data)
					if err != nil {
						panic(err)
					}
					if len(link.Length) == 0 {
						link.Length = res.Header.Get("Content-Length")
					}
					if len(link.Type) == 0 {
						link.Type = res.Header.Get("Content-Type")
					}

					link.Data = nil
					entry.Link[linkIndex] = link
				}
				entry.Link[linkIndex] = link.SetHrefLang(*f.Lang)
			}
		}

		// reset predefined
		f.Self = nil
		f.Describedby = nil
		f.Search = nil
		f.Up = nil

		processedFeeds = append(processedFeeds, f)
	}
	return processedFeeds
}
