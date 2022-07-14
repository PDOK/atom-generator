package feeds

import (
	"errors"
	"os"
	"strings"
	"testing"
)

func TestGenerateATOM(t *testing.T) {
	var updated = "2012-03-31T13:45:03Z"
	var recentupdated = "2021-10-01T00:00:00Z"
	var tests = []struct {
		input    Feeds
		updated  *string
		expected string
	}{
		0: {input: Feeds{Feeds: []Feed{Feed{InspireDls: "http://inspire.ec.europa.eu/schemas/inspire_dls/1.0",
			Lang:     sp("en"),
			ID:       "http://xyz.org/download/en.xml",
			Title:    "XYZ Example INSPIRE Download Service",
			Subtitle: "INSPIRE Download Service of organisation XYZ providing Hydrography data",
			Self: &Link{
				Href:  "http://xyz.org/download/en.xml",
				Type:  "application/atom+xml",
				Title: "This document",
			},
			Describedby: &Link{
				Href: "http://xyz.org/metadata/iso19139_document.xml",
				Type: "application/xml",
			},
			Link: []Link{
				{
					Rel:   "search",
					Href:  "http://xyz.org/search/opensearchdescription.xml",
					Type:  "application/opensearchdescription+xml",
					Title: "Open Search Description for XYZ download service",
				},
				{
					Href:     "http://xyz.org/download/de.xml",
					Rel:      "alternate",
					Type:     "application/atom+xml",
					Hreflang: sp("de"),
					Title:    "The download service information in German",
				},
				{
					Href:  "http://xyz.org/download/index.html",
					Rel:   "alternate",
					Type:  "text/html",
					Title: "An HTML version of this document",
				},
				{
					Href:     "http://xyz.org/download/index.de.html",
					Rel:      "alternate",
					Type:     "text/html",
					Hreflang: sp("de"),
					Title:    "An HTML version of this document in German",
				},
			},
			Rights:  "Copyright (c) 2012, XYZ; all rights reserved",
			Updated: &updated,
			Author: Author{
				Name:  "John Doe",
				Email: "doe@xyz.org",
			},
			Entry: []Entry{
				{
					ID:                                "http://xyz.org/data/waternetwork_feed.xml",
					Rights:                            "Copyright (c) 2002-2011, XYZ; all rights reserved",
					Updated:                           &updated,
					Summary:                           "This is the entry for water network ABC Dataset",
					Polygon:                           "47.202 5.755 55.183 5.755 55.183 15.253 47.202 15.253 47.202 5.755",
					Title:                             "Water network ABC Dataset Feed",
					SpatialDatasetIdentifierCode:      "wn_id1",
					SpatialDatasetIdentifierNamespace: "http://xyz.org/",
					Link: []Link{
						{
							Rel:  "describedby",
							Href: "http://xyz.org/metadata/abcISO19139.xml",
							Type: "application/xml",
						},
						{
							Rel:   "alternate",
							Href:  "http://xyz.org/data/waternetwork_feed.xml",
							Type:  "application/atom+xml",
							Title: "Feed containing the pre-defined waternetwork dataset (in one or more downloadable formats)",
						},
						{
							Rel:   "related",
							Href:  "http://xyz.org/wfs?request=GetCapabilities&service=WFS&version=2.0.0",
							Type:  "application/xml",
							Title: "Service implementing Direct Access operations",
						},
					},
					Category: []Category{
						{
							Term:  "http://www.opengis.net/def/crs/EPSG/0/25832",
							Label: "ETRS89 / UTM zone 32N",
						},
						{
							Term:  "http://www.opengis.net/def/crs/EPSG/0/4258",
							Label: "ETRS89",
						},
					},
				},
			},
		},
		},
		},
			expected: `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom" xmlns:georss="http://www.georss.org/georss" xmlns:inspire_dls="http://inspire.ec.europa.eu/schemas/inspire_dls/1.0" xml:lang="en">
 <id>http://xyz.org/download/en.xml</id>
 <title>XYZ Example INSPIRE Download Service</title>
 <subtitle>INSPIRE Download Service of organisation XYZ providing Hydrography data</subtitle>
 <link href="http://xyz.org/search/opensearchdescription.xml" rel="search" type="application/opensearchdescription+xml" hreflang="en" title="Open Search Description for XYZ download service"></link>
 <link href="http://xyz.org/download/de.xml" rel="alternate" type="application/atom+xml" hreflang="de" title="The download service information in German"></link>
 <link href="http://xyz.org/download/index.html" rel="alternate" type="text/html" hreflang="en" title="An HTML version of this document"></link>
 <link href="http://xyz.org/download/index.de.html" rel="alternate" type="text/html" hreflang="de" title="An HTML version of this document in German"></link>
 <link href="http://xyz.org/download/en.xml" rel="self" type="application/atom+xml" hreflang="en" title="This document"></link>
 <link href="http://xyz.org/metadata/iso19139_document.xml" rel="describedby" type="application/xml" hreflang="en"></link>
 <rights>Copyright (c) 2012, XYZ; all rights reserved</rights>
 <updated>2012-03-31T13:45:03Z</updated>
 <author>
  <name>John Doe</name>
  <email>doe@xyz.org</email>
 </author>
 <entry>
  <id>http://xyz.org/data/waternetwork_feed.xml</id>
  <title>Water network ABC Dataset Feed</title>
  <summary>This is the entry for water network ABC Dataset</summary>
  <link href="http://xyz.org/metadata/abcISO19139.xml" rel="describedby" type="application/xml" hreflang="en"></link>
  <link href="http://xyz.org/data/waternetwork_feed.xml" rel="alternate" type="application/atom+xml" hreflang="en" title="Feed containing the pre-defined waternetwork dataset (in one or more downloadable formats)"></link>
  <link href="http://xyz.org/wfs?request=GetCapabilities&amp;service=WFS&amp;version=2.0.0" rel="related" type="application/xml" hreflang="en" title="Service implementing Direct Access operations"></link>
  <rights>Copyright (c) 2002-2011, XYZ; all rights reserved</rights>
  <updated>2012-03-31T13:45:03Z</updated>
  <georss:polygon>47.202 5.755 55.183 5.755 55.183 15.253 47.202 15.253 47.202 5.755</georss:polygon>
  <category term="http://www.opengis.net/def/crs/EPSG/0/25832" label="ETRS89 / UTM zone 32N"></category>
  <category term="http://www.opengis.net/def/crs/EPSG/0/4258" label="ETRS89"></category>
  <inspire_dls:spatial_dataset_identifier_code>wn_id1</inspire_dls:spatial_dataset_identifier_code>
  <inspire_dls:spatial_dataset_identifier_namespace>http://xyz.org/</inspire_dls:spatial_dataset_identifier_namespace>
 </entry>
</feed>`},
		1: {input: Feeds{Feeds: []Feed{Feed{InspireDls: "http://inspire.ec.europa.eu/schemas/inspire_dls/1.0",
			Lang:     sp("nl"),
			ID:       "https://service.pdok.nl/kadaster/plu/atom/v1_0/plu.xml",
			Title:    "INSPIRE Download Service van Ruimtelijke plannen",
			Subtitle: "Voorgedefinieerde dataset INSPIRE download service",
			Link: []Link{
				{
					Rel:  "self",
					Href: "https://service.pdok.nl/kadaster/plu/atom/v1_0/plu.xml",
				},
				{
					Rel:   "up",
					Href:  "https://service.pdok.nl/kadaster/plu/atom/v1_0/index.xml",
					Type:  "application/atom+xml",
					Title: "Top Atom Download Service Feed",
				},
				{
					Rel:  "describedby",
					Href: "https://www.nationaalgeoregister.nl/geonetwork/srv/dut/catalog.search#/metadata/17716ed7-ce0d-4bfd-8868-a398e5578a36",
					Type: "text/html",
				},
				{
					Rel:   "related",
					Href:  "https://www.nationaalgeoregister.nl/geonetwork/srv/dut/catalog.search#/metadata/17716ed7-ce0d-4bfd-8868-a398e5578a36",
					Type:  "text/html",
					Title: "NGR pagina voor deze dataset",
				},
			},
			Rights:  "http://creativecommons.org/publicdomain/zero/1.0/deed.nl",
			Updated: &recentupdated,
			Author: Author{
				Name:  "PDOK Beheer",
				Email: "beheerPDOK@kadaster.nl",
			},
			Entry: []Entry{
				{
					ID:      "https://service.pdok.nl/kadaster/plu/atom/v1_0/plu.xml",
					Rights:  "http://creativecommons.org/publicdomain/zero/1.0/deed.nl",
					Updated: &recentupdated,
					Polygon: "50.6 3.1 50.6 7.3 53.7 7.3 53.7 3.1 50.6 3.1",
					Title:   "INSPIRE Download Service van Ruimtelijke plannen",
					Content: "Bestand is opgesplitst per featuretype, elk featuretype heeft een eigen download bestand",
					Link: []Link{
						{
							Rel:    "section",
							Href:   "https://service.pdok.nl/kadaster/plu/atom/v1_0/downloads/Besluitgebied_A.gml.gz",
							Type:   "application/x-gmz",
							Length: "3547244",
						},
						{
							Rel:    "section",
							Href:   "https://service.pdok.nl/kadaster/plu/atom/v1_0/downloads/Besluitgebied_P.gml.gz",
							Type:   "application/x-gmz",
							Length: "15714976",
						},
						{
							Rel:    "section",
							Href:   "https://service.pdok.nl/kadaster/plu/atom/v1_0/downloads/Besluitgebied_X.gml.gz",
							Type:   "application/x-gmz",
							Length: "45621084",
						},
					},
					Category: []Category{
						{
							Term:  "https://www.opengis.net/def/crs/EPSG/0/28992",
							Label: "Amersfoort / RD New",
						},
					},
				},
			},
		},
		}},
			expected: `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom" xmlns:georss="http://www.georss.org/georss" xmlns:inspire_dls="http://inspire.ec.europa.eu/schemas/inspire_dls/1.0" xml:lang="nl">
 <id>https://service.pdok.nl/kadaster/plu/atom/v1_0/plu.xml</id>
 <title>INSPIRE Download Service van Ruimtelijke plannen</title>
 <subtitle>Voorgedefinieerde dataset INSPIRE download service</subtitle>
 <link href="https://service.pdok.nl/kadaster/plu/atom/v1_0/plu.xml" rel="self" hreflang="nl"></link>
 <link href="https://service.pdok.nl/kadaster/plu/atom/v1_0/index.xml" rel="up" type="application/atom+xml" hreflang="nl" title="Top Atom Download Service Feed"></link>
 <link href="https://www.nationaalgeoregister.nl/geonetwork/srv/dut/catalog.search#/metadata/17716ed7-ce0d-4bfd-8868-a398e5578a36" rel="describedby" type="text/html" hreflang="nl"></link>
 <link href="https://www.nationaalgeoregister.nl/geonetwork/srv/dut/catalog.search#/metadata/17716ed7-ce0d-4bfd-8868-a398e5578a36" rel="related" type="text/html" hreflang="nl" title="NGR pagina voor deze dataset"></link>
 <rights>http://creativecommons.org/publicdomain/zero/1.0/deed.nl</rights>
 <updated>2021-10-01T00:00:00Z</updated>
 <author>
  <name>PDOK Beheer</name>
  <email>beheerPDOK@kadaster.nl</email>
 </author>
 <entry>
  <id>https://service.pdok.nl/kadaster/plu/atom/v1_0/plu.xml</id>
  <title>INSPIRE Download Service van Ruimtelijke plannen</title>
  <content>Bestand is opgesplitst per featuretype, elk featuretype heeft een eigen download bestand</content>
  <link href="https://service.pdok.nl/kadaster/plu/atom/v1_0/downloads/Besluitgebied_A.gml.gz" rel="section" type="application/x-gmz" hreflang="nl" length="3547244"></link>
  <link href="https://service.pdok.nl/kadaster/plu/atom/v1_0/downloads/Besluitgebied_P.gml.gz" rel="section" type="application/x-gmz" hreflang="nl" length="15714976"></link>
  <link href="https://service.pdok.nl/kadaster/plu/atom/v1_0/downloads/Besluitgebied_X.gml.gz" rel="section" type="application/x-gmz" hreflang="nl" length="45621084"></link>
  <rights>http://creativecommons.org/publicdomain/zero/1.0/deed.nl</rights>
  <updated>2021-10-01T00:00:00Z</updated>
  <georss:polygon>50.6 3.1 50.6 7.3 53.7 7.3 53.7 3.1 50.6 3.1</georss:polygon>
  <category term="https://www.opengis.net/def/crs/EPSG/0/28992" label="Amersfoort / RD New"></category>
 </entry>
</feed>`},
	}

	for k, test := range tests {
		p := ProcessFeeds(test.input)
		output := p[0].GenerateATOM()
		if string(output) != test.expected {
			t.Errorf("test: %d, expected: %s \ngot: %s", k, test.expected, string(output))
		}
	}
}

func TestGetFileName(t *testing.T) {
	var tests = []struct {
		input    Feed
		expected string
	}{
		0: {input: Feed{ID: `http://xyz.org/download/en.xml`}, expected: `en.xml`},
		1: {input: Feed{ID: `not a URL.xml`}, expected: "Not a valid ID was provided, got: `not a URL.xml`"},
	}

	for k, test := range tests {
		output, err := test.input.GetFileName()
		if err == nil {
			if output != test.expected {
				t.Errorf("test: %d, expected: %s \ngot: %s", k, test.expected, output)
			}
		} else {
			if err.Error() != test.expected {
				t.Errorf("test: %d, expected: %s \ngot: %s", k, test.expected, err)
			}
		}
	}
}

func TestValid(t *testing.T) {
	var updated = "2021-03-31T13:45:03Z"
	var tests = []struct {
		input    Feed
		expected error
	}{
		0: {
			input: Feed{
				ID:      "http://xyz.org/download/en.xml",
				Rights:  "Copyright (c) 2012, XYZ; all rights reserved",
				Updated: &updated,
				Author: Author{
					Name:  "John Doe",
					Email: "doe@xyz.org",
				},
			},
			expected: nil,
		},
		1: {
			input: Feed{
				ID:      "http://xyz.org/download/en.xml",
				Rights:  "Copyright (c) 2012, XYZ; all rights reserved",
				Updated: &updated,
			},
			expected: errors.New(invalidauthor),
		},
		2: {
			input: Feed{
				ID:      "http://xyz.org/download/en.xml",
				Updated: &updated,
				Author: Author{
					Name:  "John Doe",
					Email: "doe@xyz.org",
				},
			},
			expected: errors.New(invalidrights),
		},
		3: {
			input: Feed{
				ID:     "http://xyz.org/download/en.xml",
				Rights: "Copyright (c) 2012, XYZ; all rights reserved",
				Author: Author{
					Name:  "John Doe",
					Email: "doe@xyz.org",
				},
			},
			expected: errors.New(invaliddatetime),
		},
		4: {
			input: Feed{
				ID: "xyzorgdownloaden.xml",
			},
			expected: errors.New(invalidid),
		},
	}

	for k, test := range tests {
		b := test.input.Valid()
		if b == nil {
			if b != test.expected {
				t.Errorf("test: %d, expected: %t \ngot: %t", k, test.expected, b)
			}
		} else {
			if test.expected != nil {
				if b.Error() != test.expected.Error() {
					t.Errorf("test: %d, expected: %t \ngot: %t", k, test.expected, b)
				}
			} else {
				t.Errorf("test: %d, expected: %t \ngot: %t", k, test.expected, b)
			}
		}
	}
}

func TestFeedWriteATOM(t *testing.T) {
	type fields struct {
	}
	type args struct {
		filename string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		0: {
			fields: fields{},
			args: args{
				"test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Feed{}

			f.WriteATOM(tt.args.filename)
			_, err := os.Stat(tt.args.filename)
			err2 := os.Remove(tt.args.filename)
			if err2 != nil {
				t.Errorf("Error occurred: %s", err2)
			}
			if os.IsNotExist(err) {
				t.Errorf("File not created.")
			}
		})
	}
}

func TestFeedStyleSheet(t *testing.T) {

	var testEmpty = ""
	var test = "test"

	var tests = []struct {
		input    Feed
		expected string
	}{
		0: {input: Feed{XMLStylesheet: &testEmpty}, expected: ""},
		1: {input: Feed{XMLStylesheet: &test}, expected: `<?xml-stylesheet href="test" type="text/xsl" media="screen"?>`},
	}
	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			got := tt.input.StyleSheet()
			if !strings.Contains(string(got), tt.expected) {
				t.Errorf("StyleSheet() = %v, want %v", string(got), tt.expected)
			}
		})
	}
}
