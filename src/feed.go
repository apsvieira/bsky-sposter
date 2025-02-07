package sposter

import (
	"fmt"
	"time"

	"github.com/mmcdole/gofeed"
)

func FetchNewItems(feedURL string, lastParsed *time.Time) ([]*gofeed.Item, error) {
	// Create a new parser
	fp := gofeed.NewParser()

	// Parse the feed
	feed, err := fp.ParseURL(feedURL)
	if err != nil {
		return nil, fmt.Errorf("error fetching feed: %v", err)
	}

	if lastParsed == nil {
		return feed.Items, nil
	}

	filtered := make([]*gofeed.Item, 0)
	for _, item := range feed.Items {
		if item.PublishedParsed.After(*lastParsed) {
			filtered = append(filtered, item)
		}
	}

	return filtered, nil
}
