package feeds

import (
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

// Feeds struct
type Feeds struct {
	Feeds []Feed `yaml:"feeds"`
}

// Feed struct
type Feed struct {
	XMLName       xml.Name `xml:"feed"`
	XMLStylesheet *string  `yaml:"stylesheet"`
	Xmlns         string   `xml:"xmlns,attr" yaml:"xmlns"`                                       //"http://www.w3.org/2005/Atom"
	Georss        string   `xml:"xmlns:georss,attr,omitempty" yaml:"georss,omitempty"`           //"http://www.georss.org/georss"
	InspireDls    string   `xml:"xmlns:inspire_dls,attr,omitempty" yaml:"inspire_dls,omitempty"` //"http://inspire.ec.europa.eu/schemas/inspire_dls/1.0"
	Lang          *string  `xml:"xml:lang,attr,omitempty" yaml:"lang,omitempty"`

	ID       string `xml:"id" yaml:"id"`
	Title    string `xml:"title" yaml:"title"`
	Subtitle string `xml:"subtitle" yaml:"subtitle"`

	// Placeholder Links, need to be moved to []Link and deleted
	Self        *Link `xml:"self,omitempty" yaml:"self,omitempty"`
	Describedby *Link `xml:"describedby,omitempty" yaml:"describedby,omitempty"`
	Search      *Link `xml:"search,omitempty" yaml:"search,omitempty"`
	Up          *Link `xml:"up,omitempty" yaml:"up,omitempty"`

	Link []Link `xml:"link" yaml:"link"`

	Rights  string  `xml:"rights" yaml:"rights"`
	Updated *string `xml:"updated" yaml:"updated,omitempty"`
	Author  Author  `xml:"author" yaml:"author"`
	Entry   []Entry `xml:"entry" yaml:"entry"`
}

// GetFileName function
// extracts a filename based on the ID element of the Feed struct
// TG Requirement 9 - Technical Guidance Download Services v3.1
// The 'id' element of a feed shall contain an HTTP URI which dereferences to the feed
func (f *Feed) GetFileName() (string, error) {
	if !strings.Contains(f.ID, `http`) {
		return ``, fmt.Errorf("not a valid ID was provided, got: `%s`", f.ID)
	}

	parts := strings.Split(f.ID, `/`)
	filename := parts[len(parts)-1]

	return filename, nil
}

// GenerateATOM function build a ATOM feed from the configuration
func (f *Feed) GenerateATOM() []byte {
	stylesheet := f.StyleSheet()
	f.XMLStylesheet = nil

	si, _ := xml.MarshalIndent(f, "", " ")
	return append(append([]byte(xml.Header), stylesheet...), si...)
}

// WriteATOM function writes the ATOM feed to file
func (f *Feed) WriteATOM(filename string) {
	b := f.GenerateATOM()
	err := os.WriteFile(filename, b, 0777)
	if err != nil {
		log.Fatalf("Could not write to file %s : %v ", filename, err)
	}
}

// StyleSheet function returns a xml-stylesheet header if available
func (f *Feed) StyleSheet() []byte {
	if f.XMLStylesheet != nil {
		header := `<?xml-stylesheet href="` + *f.XMLStylesheet + `" type="text/xsl" media="screen"?>` + "\n"
		return []byte(header)
	}

	return []byte(``)
}

// Valid function that validates the Feed based on TG Requirements
// For now a simple validation
func (f *Feed) Valid() error {

	// TG Requirement 5
	// The 'title' element of an Atom feed shall be populated with a human readable title for the feed.
	if len(f.Title) == 0 {
		return errors.New(invalidtitle)
	}

	// TG Recommendation 1
	// The 'subtitle' element of an Atom feed may be populated with a human readable subtitle for the feed.
	if len(f.Subtitle) == 0 {
		log.Println(warningsubtitle)
	}

	// TG Requirement 9
	// The 'id' element of a feed shall contain an HTTP URI which dereferences to the feed
	_, err := url.ParseRequestURI(f.ID)
	if err != nil {
		return errors.New(invalidid)
	}

	// TG Requirement 10
	// The 'rights' element of a feed shall contain information about rights or restrictions for that feed.
	if len(f.Rights) == 0 {
		return errors.New(invalidrights)
	}

	// TG Requirement 11
	// The 'updated' element of a feed shall contain the date, time and timezone at which the feed was last updated.
	for _, entry := range f.Entry {
		if entry.Updated == nil {
			return errors.New(invalidupdated)
		}
		if _, err := time.Parse(`2006-01-02T15:04:05Z`, *entry.Updated); err != nil {
			return errors.New(invaliddatetime)
		}
	}
	if f.Updated == nil {
		return errors.New(invalidupdated)
	}
	if _, err := time.Parse(`2006-01-02T15:04:05Z`, *f.Updated); err != nil {
		return errors.New(invaliddatetime)
	}

	// TG Recommendation 11
	// Where a dataset is provided in multiple physical files: a `time` attribute may be used to describe the temporal extent of a particular file.
	// If this is used, then the value of this attribute should be structured according to the ISO 8601 standard.
	for _, entry := range f.Entry {
		for _, link := range entry.Link {
			if link.Time == nil {
				continue
			}
			if _, err := time.Parse(`2006-01-02T15:04:05Z`, *link.Time); err != nil {
				return errors.New(invalidlinktime)
			}
		}
	}

	// TG Recommendation 10
	// Where a dataset is provided in multiple physical files: a `bbox` attribute may be used to describe the geospatial extent of a particular file.
	// If this is used, then the value of this attribute should be structured according to the georss:box structure.
	for _, entry := range f.Entry {
		for _, link := range entry.Link {
			if link.Bbox == nil {
				continue
			}
			matched, err := regexp.MatchString(`^-?\d+(\.\d+)? -?\d+(\.\d+)? -?\d+(\.\d+)? -?\d+(\.\d+)?$`, *link.Bbox)
			if !matched {
				return errors.New(invalidlinkbbox)
			}
			if err != nil {
				return errors.New("error while matching bbox regex: " + err.Error())
			}
		}
	}

	// TG Requirement 12
	// The 'author' element of a feed shall contain current contact information for an individual or organisation responsible for the feed. At the minimum, a name and email address shall be provided as contact information.
	if len(f.Author.Name) == 0 || len(f.Author.Email) == 0 {
		return errors.New(invalidauthor)
	}

	return nil
}

// Function that retrieves values from updated fields and returns the most recent updated field
func (f *Feed) recentUpdated(feeds Feeds) {
	for index, entry := range f.Entry {
		if entry.Updated == nil {
			nestedfeed := entry.nestedFeed(feeds.Feeds)
			if nestedfeed != nil {
				updated := nestedfeed.recentUpdatedEntry()
				f.Entry[index].Updated = updated
			}
		}
	}
	if f.Updated == nil {
		feedUpdated := f.recentUpdatedEntry()
		f.Updated = feedUpdated
	}
}

func (e Entry) nestedFeed(feeds []Feed) *Feed {
	for _, feed := range feeds {
		if e.ID == feed.ID {
			return &feed
		}
	}
	return nil
}

func (f *Feed) recentUpdatedEntry() *string {
	var updatedfields []string

	for _, entry := range f.Entry {
		if entry.Updated == nil {
		} else {
			updatedfields = append(updatedfields, *entry.Updated)
		}
	}
	sort.Sort(sort.Reverse(sort.StringSlice(updatedfields)))
	if len(updatedfields) == 0 {
		return nil
	} else {
		return &updatedfields[0]
	}

}

// Entry struct
type Entry struct {
	ID                                string     `xml:"id" yaml:"id"`
	Title                             string     `xml:"title,omitempty" yaml:"title,omitempty"`
	Content                           string     `xml:"content,omitempty" yaml:"content,omitempty"`
	Summary                           string     `xml:"summary,omitempty" yaml:"summary,omitempty"`
	Link                              []Link     `xml:"link" yaml:"link"`
	Rights                            string     `xml:"rights,omitempty" yaml:"rights,omitempty"`
	Updated                           *string    `xml:"updated" yaml:"updated,omitempty"`
	Polygon                           string     `xml:"georss:polygon,omitempty" yaml:"polygon,omitempty"`
	Category                          []Category `xml:"category" yaml:"category"`
	SpatialDatasetIdentifierCode      *string    `xml:"inspire_dls:spatial_dataset_identifier_code,omitempty" yaml:"spatial_dataset_identifier_code,omitempty"`
	SpatialDatasetIdentifierNamespace *string    `xml:"inspire_dls:spatial_dataset_identifier_namespace,omitempty" yaml:"spatial_dataset_identifier_namespace,omitempty"`
}

// Author struct
type Author struct {
	Name  string `xml:"name" yaml:"name"`
	Email string `xml:"email" yaml:"email"`
}

// Link struct
type Link struct {
	Href     string  `xml:"href,attr" yaml:"href"`
	Data     *string `yaml:"data"`
	Rel      string  `xml:"rel,attr,omitempty" yaml:"rel,omitempty"`
	Type     string  `xml:"type,attr,omitempty" yaml:"type,omitempty"`
	Hreflang *string `xml:"hreflang,attr,omitempty" yaml:"hreflang,omitempty"`
	Length   string  `xml:"length,attr,omitempty" yaml:"length,omitempty"`
	Title    string  `xml:"title,attr,omitempty" yaml:"title,omitempty"`
	Version  *string `xml:"version,attr,omitempty" yaml:"version,omitempty"`
	Time     *string `xml:"time,attr,omitempty" yaml:"time,omitempty"`
	Bbox     *string `xml:"bbox,attr,omitempty" yaml:"bbox,omitempty"`
}

// SetHrefLang function assigns a default Lang is none is given
func (l *Link) SetHrefLang(lang string) Link {
	if l.Hreflang == nil {
		l.Hreflang = &lang
	}

	return *l
}

// Category struct
type Category struct {
	Term  string `xml:"term,attr" yaml:"term"`
	Label string `xml:"label,attr" yaml:"label"`
}
