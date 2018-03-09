package feeds

type Feeds struct {
	Feeds []Feed `json:"feeds"`
}

type Feed struct {
	URL      string `json:"url"`
	Articles []interface{} `json:"articles"`
}

func New() Feeds {
	return Feeds{[]Feed{}}
}
