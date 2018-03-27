package pocket

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

type User struct {
	AccessToken string
	ConsumerKey string
}

type List struct {
	Status int             `json:"status"`
	Error  string          `json:"error"`
	List   map[string]Item `json:"list"`
}

type Item struct {
	ID            string `json:"item_id"`
	ResolvedURL   string `json:"resolved_url"`
	ResolvedTitle string `json:"resolved_title"`
	Excerpt       string `json:"excerpt"`
}

var endpoint string = "https://getpocket.com/v3/get"

func NewUser(token, key string) User {
	return User{
		AccessToken: token,
		ConsumerKey: key,
	}
}

func (u User) GetList() (list List, err error) {
	requestParams, err := json.Marshal(map[string]string{
		"consumer_key": u.ConsumerKey,
		"access_token": u.AccessToken,
		"detailType":   "simple",
	})
	if err != nil {
		return list, err
	}

	resp, err := http.Post(
		endpoint,
		"application/json",
		strings.NewReader(string(requestParams)),
	)
	if err != nil {
		return list, err
	}

	if resp.StatusCode != 200 {
		return list, errors.New(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return list, err
	}

	err = json.Unmarshal(body, &list)
	if err != nil {
		return list, err
	}

	if list.Error != "" {
		return list, errors.New(list.Error)
	}

	return list, nil
}

func (u User) GetUnreadCount() (int, error) {
	list, err := u.GetList()
	if err != nil {
		return 0, err
	}

	return len(list.List), nil
}
