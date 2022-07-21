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

const (
	invalidtitle    = "invalid 'title' need to be human readable title for the feed see TG Requirement 5"
	invalidid       = "invalid 'id', needs to be an HTTP URI see TG Requirement 9"
	invalidrights   = "invalid 'rights', cannot be empty see TG Requirement 10"
	invaliddatetime = "invalid 'updated', needs to be a valid datetime with timezone see TG Requirement 11"
	invalidauthor   = "invalid 'author', cannot be empty see TG Requirement 12"
	invalidupdated  = "invalid 'updated', updated is required see TG Requirements 11"
)

const (
	warningsubtitle = "missing 'subtitle' may be a human readable subtitle for the feed see TG Recommandation 1"
)

// GetDefaultFeedProperties returns mandatory/static ServiceFeed properties
func GetDefaultFeedProperties() Feed {
	var f Feed

	f.Xmlns = "http://www.w3.org/2005/Atom"   // TG Requirement 2 - Technical Guidance Download Services v3.1
	f.Georss = "http://www.georss.org/georss" // TG Requirement 3 - Technical Guidance Download Services v3.1

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
