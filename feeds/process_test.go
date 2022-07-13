package feeds

import (
	"reflect"
	"testing"
)

func sp(s string) *string {
	return &s
}

func TestProcessFeed(t *testing.T) {
	var updated = "2021-06-15T11:12:34Z"
	var tests = []struct {
		input    Feeds
		expected []Feed
	}{
		0: {input: Feeds{Feeds: []Feed{Feed{
			ID:    "http://xyz.org/download/en.xml",
			Title: "XYZ Example INSPIRE Download Service",
			Self: &Link{
				Href:     "http://xyz.org/download/en.xml",
				Type:     "application/atom+xml",
				Hreflang: sp("en"),
				Title:    "This document"},
			Link: []Link{
				{
					Href:  "http://xyz.org/search/opensearchdescription.xml",
					Rel:   "search",
					Type:  "application/opensearchdescription+xml",
					Title: "Open Search Description for XYZ download service",
				},
			},
		}}},
			expected: []Feed{Feed{
				Xmlns: "http://www.w3.org/2005/Atom", Georss: "http://www.georss.org/georss", Lang: sp("en"),
				ID:    "http://xyz.org/download/en.xml",
				Title: "XYZ Example INSPIRE Download Service",
				Link: []Link{
					{
						Href:     "http://xyz.org/search/opensearchdescription.xml",
						Rel:      "search",
						Type:     "application/opensearchdescription+xml",
						Title:    "Open Search Description for XYZ download service",
						Hreflang: sp("en"),
					},
					{
						Href:     "http://xyz.org/download/en.xml",
						Rel:      "self",
						Type:     "application/atom+xml",
						Hreflang: sp("en"),
						Title:    "This document",
					},
				},
			}},
		},
		1: {input: Feeds{Feeds: []Feed{Feed{
			ID:    "http://xyz.org/download/en.xml",
			Title: "XYZ Example INSPIRE Download Service",
			Self: &Link{
				Href:     "http://xyz.org/download/en.xml",
				Type:     "application/atom+xml",
				Hreflang: sp("en"),
				Title:    "This document"},
			Describedby: &Link{
				Href: "http://xyz.org/metadata/iso19139_document.xml",
				Type: "application/xml"},
			Search: &Link{
				Href:  "http://xyz.org/search/opensearchdescription.xml",
				Type:  "application/opensearchdescription+xml",
				Title: "Open Search Description for XYZ download service"},
			Up: &Link{
				Href:     "http://xyz.org/download/en.xml",
				Type:     "application/atom+xml",
				Hreflang: sp("en"),
				Title:    "This document"},
			Link: []Link{
				{
					Rel:      "alternate",
					Href:     "http://xyz.org/download/index.de.html",
					Type:     "text/html",
					Hreflang: sp("de"),
					Title:    "An HTML version of this document in German",
				},
			},
		}}},
			expected: []Feed{Feed{
				Xmlns: "http://www.w3.org/2005/Atom", Georss: "http://www.georss.org/georss", Lang: sp("en"),
				ID:    "http://xyz.org/download/en.xml",
				Title: "XYZ Example INSPIRE Download Service",
				Link: []Link{
					{
						Rel:      "alternate",
						Href:     "http://xyz.org/download/index.de.html",
						Type:     "text/html",
						Hreflang: sp("de"),
						Title:    "An HTML version of this document in German",
					},
					{
						Rel:      "self",
						Href:     "http://xyz.org/download/en.xml",
						Type:     "application/atom+xml",
						Hreflang: sp("en"),
						Title:    "This document",
					},
					{
						Rel:      "describedby",
						Href:     "http://xyz.org/metadata/iso19139_document.xml",
						Type:     "application/xml",
						Hreflang: sp("en"),
					},
					{
						Rel:      "search",
						Href:     "http://xyz.org/search/opensearchdescription.xml",
						Type:     "application/opensearchdescription+xml",
						Title:    "Open Search Description for XYZ download service",
						Hreflang: sp("en"),
					},
					{
						Rel:      "up",
						Href:     "http://xyz.org/download/en.xml",
						Type:     "application/atom+xml",
						Hreflang: sp("en"),
						Title:    "This document",
					},
				},
			}},
		},
		2: {input: Feeds{Feeds: []Feed{Feed{
			ID:    "http://xyz.org/download/en.xml",
			Title: "Service Feed",
			Entry: []Entry{
				Entry{
					ID: "datafeed-1",
				},
			},
		},
			Feed{
				ID:    "datafeed-1",
				Title: "Data Feed",
				Entry: []Entry{
					Entry{
						ID:      "download-entry",
						Updated: &updated,
					},
				},
			},
		}},
			expected: []Feed{Feed{
				Xmlns: "http://www.w3.org/2005/Atom", Georss: "http://www.georss.org/georss", Lang: sp("en"),
				ID:      "http://xyz.org/download/en.xml",
				Title:   "Service Feed",
				Updated: &updated,
				Entry: []Entry{
					Entry{
						ID:      "datafeed-1",
						Updated: &updated,
					},
				},
			},
				Feed{
					Xmlns: "http://www.w3.org/2005/Atom", Georss: "http://www.georss.org/georss", Lang: sp("en"),
					ID:      "datafeed-1",
					Title:   "Data Feed",
					Updated: &updated,
					Entry: []Entry{
						Entry{
							ID:      "download-entry",
							Updated: &updated,
						},
					},
				},
			},
		},
	}

	for k, test := range tests {
		output := ProcessFeeds(test.input)
		if !reflect.DeepEqual(output, test.expected) {
			t.Errorf("test: %d, expected: \n%#v+ \ngot: \n%#v+", k, test.expected, output)
		}
	}
}
