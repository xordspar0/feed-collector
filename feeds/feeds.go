package feeds

type Feeds struct {
	Feeds []Feed `json:"feeds"`
}

type Feed struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Count int    `json:"count"`
}

func New() Feeds {
	return Feeds{[]Feed{}}
}

func (f *Feeds) AddFeed(name, url string) *Feed {
	f.Feeds = append(f.Feeds, Feed{
		Name: name,
		URL:  url,
	})
	return &f.Feeds[len(f.Feeds)-1]
}
