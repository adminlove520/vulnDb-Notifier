package rss

import (
	"github.com/adminlove520/vulnDb-Notifier/internal/errors"
	"github.com/mmcdole/gofeed"
)

func ParseFeed(feedURL string) (*gofeed.Feed, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(feedURL)
	if err != nil {
		return nil, &errors.RSSFeedError{Message: "Failed to parse RSS feed: " + err.Error()}
	}

	return feed, nil
}
