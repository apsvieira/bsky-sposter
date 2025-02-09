package main

import (
	"context"
	"log"
	"time"

	sposter "github.com/apsvieira/bsky-sposter/src"
	"github.com/apsvieira/bsky-sposter/src/atproto/client"
)

func main() {
	ctx := context.Background()
	creds, err := client.GetCredentials()
	if err != nil {
		log.Fatalf("Error getting credentials: %s", err)
	}

	client, err := client.NewClient(ctx, sposter.BLUESKY_SERVICE, creds)
	if err != nil {
		log.Fatalf("Error authenticating: %s", err)
	}

	minDate := time.Date(2025, time.February, 1, 0, 0, 0, 0, time.UTC)
	articles, err := sposter.FetchNewItems("https://apsv.bearblog.dev/feed/", &minDate)

	if err != nil {
		log.Panicf("Error fetching feed: %s", err)
	}

	log.Printf("Found %d new posts", len(articles))
	for _, item := range articles {
		p, err := sposter.NewPostFromFeedItem(item)
		if err != nil {
			log.Fatalf("Error creating post: %s", err)
			continue
		}

		msg, err := p.BskyPost()
		if err != nil {
			log.Fatalf("Error creating message: %s", err)
			continue
		}

		post, err := sposter.NewPost(ctx, client, msg)
		if err != nil {
			log.Fatalf("Error creating post: %s", err)
		}

		bskyPost, err := client.CreatePost(ctx, post)
		if err != nil {
			log.Fatalf("Error posting: %s", err)
		}
		log.Printf("Posted: %s", bskyPost.Uri)
	}
}
