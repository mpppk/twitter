package twitter

import (
	"fmt"
	"net/url"

	"github.com/ChimeraCoder/anaconda"
)

type Client struct {
	client *anaconda.TwitterApi
	maxId  int64
}

func NewClient(accessToken, accessTokenSecret, consumerKey, consumerSecret string) *Client {
	api := anaconda.NewTwitterApiWithCredentials(
		accessToken, accessTokenSecret, consumerKey, consumerSecret,
	)
	return &Client{
		client: api,
	}
}

func (c *Client) SearchTweets(query string, maxId, sinceId int64) ([]anaconda.Tweet, error) {
	v := url.Values{}
	v.Set("count", "100")
	if maxId > 0 {
		lastTweetIdStr := fmt.Sprintf("%d", maxId-1)
		v.Set("max_id", lastTweetIdStr)
	}
	if sinceId > 0 {
		sinceIdStr := fmt.Sprintf("%d", sinceId)
		v.Set("since_id", sinceIdStr)
	}
	searchResult, err := c.client.GetSearch(query, v)
	if err != nil {
		return nil, err
	}

	return searchResult.Statuses, nil
}
