package nextcloudnews

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Host struct {
	Host     string
	User     string
	Password string
}

type Feeds struct {
	StarredCount int    `json:"starredCount"`
	Feeds        []Feed `json:"feeds"`
	NewestItemId int    `json:"newestItemId"`
}

type Feed struct {
	ID               int    `json:"id"`
	URL              string `json:"url"`
	Title            string `json:"title"`
	FaviconLink      string `json:"faviconLink"`
	Added            int    `json:"added"`
	FolderID         int    `json:"folderId"`
	UnreadCount      int    `json:"unreadCount"`
	Ordering         int    `json:"ordering"`
	Link             string `json:"link"`
	Pinned           bool   `json:"pinned"`
	UpdateErrorCount int    `json:"updateErrorCount"`
	LastUpdateError  string `json:"lastUpdateError"`
}

var feedsEndpoint string = "/index.php/apps/news/api/v1-2/feeds"

func New(host, user, password string) Host {
	return Host{
		Host:     host,
		User:     user,
		Password: password,
	}
}

func (h Host) GetFeeds() (feeds Feeds, err error) {
	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, h.Host+feedsEndpoint, nil)
	if err != nil {
		return feeds, err
	}
	req.SetBasicAuth(h.User, h.Password)

	resp, err := client.Do(req)
	if err != nil {
		return feeds, err
	}

	if resp.StatusCode != 200 {
		return feeds, fmt.Errorf(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return feeds, err
	}

	err = json.Unmarshal(body, &feeds)
	if err != nil {
		return feeds, err
	}

	return feeds, nil
}

func (h Host) GetUnreadCount() (int, error) {
	feeds, err := h.GetFeeds()
	if err != nil {
		return 0, err
	}

	var totalUnreadCount int
	for _, feed := range feeds.Feeds {
		totalUnreadCount += feed.UnreadCount
	}

	return totalUnreadCount, nil
}
