package feeds

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"sort"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// Feeds struct
type Feeds struct {
	Feeds []Feed `yaml:"feeds"`
}

// Feed struct
type Feed struct {
	XMLName       xml.Name `xml:"feed"`
	XMLStylesheet *string  `yaml:"stylesheet"`
	Xmlns         string   `xml:"xmlns,attr" yaml:"xmlns"`                             //"http://www.w3.org/2005/Atom"
	Georss        string   `xml:"xmlns:georss,attr,omitempty" yaml:"georss"`           //"http://www.georss.org/georss"
	InspireDls    string   `xml:"xmlns:inspire_dls,attr,omitempty" yaml:"inspire_dls"` //"http://inspire.ec.europa.eu/schemas/inspire_dls/1.0"
	Lang          *string  `xml:"xml:lang,attr,omitempty" yaml:"lang"`

	ID       string `xml:"id" yaml:"id"`
	Title    string `xml:"title" yaml:"title"`
	Subtitle string `xml:"subtitle" yaml:"subtitle"`

	// Placeholder Links, need to be moved to []Link and deleted
	Self        *Link `xml:"self,omitempty" yaml:"self"`
	Describedby *Link `xml:"describedby,omitempty" yaml:"describedby"`
	Search      *Link `xml:"search,omitempty" yaml:"search"`
	Up          *Link `xml:"up,omitempty" yaml:"up"`

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
	err := ioutil.WriteFile(filename, b, 0777)
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
			return errors.New(invaliddatetime)
		}
		if _, err := time.Parse(`2006-01-02T15:04:05Z`, *entry.Updated); err != nil {
			return errors.New(invaliddatetime)
		}
	}
	if f.Updated == nil {
		return errors.New(invaliddatetime)
	}
	if _, err := time.Parse(`2006-01-02T15:04:05Z`, *f.Updated); err != nil {
		return errors.New(invaliddatetime)
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
	Title                             string     `xml:"title,omitempty" yaml:"title"`
	Content                           string     `xml:"content,omitempty" yaml:"content"`
	Summary                           string     `xml:"summary,omitempty" yaml:"summary"`
	Link                              []Link     `xml:"link" yaml:"link"`
	Rights                            string     `xml:"rights,omitempty" yaml:"rights"`
	Updated                           *string    `xml:"updated" yaml:"updated,omitempty"`
	Polygon                           string     `xml:"georss:polygon,omitempty" yaml:"polygon"`
	Category                          []Category `xml:"category" yaml:"category"`
	SpatialDatasetIdentifierCode      string     `xml:"inspire_dls:spatial_dataset_identifier_code,omitempty" yaml:"spatial_dataset_identifier_code"`
	SpatialDatasetIdentifierNamespace string     `xml:"inspire_dls:spatial_dataset_identifier_namespace,omitempty" yaml:"spatial_dataset_identifier_namespace"`
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
	Rel      string  `xml:"rel,attr,omitempty" yaml:"rel"`
	Type     string  `xml:"type,attr,omitempty" yaml:"type"`
	Hreflang *string `xml:"hreflang,attr,omitempty" yaml:"hreflang"`
	Length   string  `xml:"length,attr,omitempty" yaml:"length"`
	Title    string  `xml:"title,attr,omitempty" yaml:"title"`
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
