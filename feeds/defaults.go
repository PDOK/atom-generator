package feeds

const (
	describedby = `describedby`
	self        = `self`
	search      = `search`
	up          = `up`

	// DEFAULTLANG contains the default language of the ATOM Feed
	// The default of the DEFAULTLANG is 'en'
	// this can be overwritten through the configuration yaml
	defaultlang = `en`
)

// GetDefaultFeedProperties returns mandatory/static ServiceFeed properties
func GetDefaultFeedProperties() Feed {
	var f Feed

	f.Xmlns = "http://www.w3.org/2005/Atom"
	f.Georss = "http://www.georss.org/georss"

	if f.Lang == nil {
		l := defaultlang
		f.Lang = &l
	}

	return f
}

// DescribedBy returns a Link containing a mandatory DescribedBy element
// TG Requirement 6 - Technical Guidance Download Services v3.1
func DescribedBy(l Link) Link {
	l.Rel = describedby
	l.Type = `application/xml`

	return l
}

// Self returns a Link containing a mandatory Self element
// TG Requirement 7 - Technical Guidance Download Services v3.1
func Self(l Link) Link {
	l.Rel = self
	l.Type = `application/atom+xml`

	return l
}

// Search returns a Link containing a Search element (which is mandatory for Service Feeds)
// TG Requirement 8 - Technical Guidance Download Services v3.1
func Search(l Link) Link {
	l.Rel = search
	l.Type = `application/opensearchdescription+xml`

	return l
}

// Up returns a Link containing a Up element (which is recommended)
// TG Recommendation 9 - Technical Guidance Download Services v3.1
func Up(l Link) Link {
	l.Rel = up
	l.Type = `application/atom+xml`

	return l
}
