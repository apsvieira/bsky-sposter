package main

import (
	"log"
	"time"

	sposter "github.com/apsvieira/bsky-sposter/src"
)

func main() {
	minDate := time.Date(2025, time.February, 1, 0, 0, 0, 0, time.UTC)
	posts, err := sposter.FetchNewItems("https://apsv.bearblog.dev/feed/", &minDate)

	if err != nil {
		log.Panicf("Error fetching feed: %s", err)
	}

	log.Printf("Found %d new posts", len(posts))
	for _, post := range posts {
		p, err := sposter.NewPostFromFeedItem(post)
		if err != nil {
			log.Printf("Error creating post: %s", err)
			continue
		}

		msg, err := p.BskyPost()
		if err != nil {
			log.Printf("Error creating message: %s", err)
			continue
		}
		log.Printf("Post: %s", msg)
	}
}
