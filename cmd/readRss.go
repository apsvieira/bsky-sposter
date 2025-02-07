package main

import (
	"log"
	"time"

	sposter "github.com/apsvieira/bsky-sposter/src"
)

func main() {
	minDate := time.Date(2025, time.February, 1, 0, 0, 0, 0, time.UTC)
	posts, err := sposter.FetchNewFromFeed("https://apsv.bearblog.dev/feed/", &minDate)

	if err != nil {
		log.Panicf("Error fetching feed: %s", err)
	}

	log.Printf("Found %d new posts", len(posts))
	for _, post := range posts {
		log.Printf("Title: %s", post.Title)
		log.Printf("Link: %s", post.Link)
		log.Printf("Published: %s", post.Published)
		log.Printf("Updated: %s", post.Updated)
	}
}
