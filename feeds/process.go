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
		_ = mergo.Merge(&f, d)

		links := f.Link
		if f.Self != nil {
			links = append(links, Self(*f.Self))
		}
		if f.Describedby != nil {
			links = append(links, DescribedBy(*f.Describedby))
		}
		if f.Search != nil {
			links = append(links, Search(*f.Search))
		}
		if f.Up != nil {
			links = append(links, Up(*f.Up))
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
					defer res.Body.Close()

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
