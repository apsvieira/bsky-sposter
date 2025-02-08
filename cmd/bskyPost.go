package main

import (
	"context"
	"log"
	"time"

	atproto "github.com/apsvieira/bsky-sposter/src/atproto"
	"github.com/apsvieira/bsky-sposter/src/richtext"
	"github.com/bluesky-social/indigo/api/bsky"
)

func main() {
	ctx := context.Background()
	creds, err := atproto.GetCredentials()
	if err != nil {
		log.Fatalf("Error getting credentials: %s", err)
	}

	client, err := atproto.NewClient(ctx, "https://bsky.social", creds)
	if err != nil {
		log.Fatalf("Error authenticating: %s", err)
	}
	log.Printf("Authenticated as %s", creds.Handle)

	// Create a new RichText instance
	rt := richtext.NewRichText("Hello, @apsv.bsky.social! #bsky https://bsky.social! google.com ")
	if err := rt.DetectFacets(ctx, client); err != nil {
		log.Fatalf("Error detecting facets: %s", err)
	}

	// Post the RichText to the feed
	post := &bsky.FeedPost{
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		Text:      rt.Text(),
		Facets:    rt.Facets(),
	}

	response, err := client.App.Bsky.Feed.Post.Create(ctx, post)
	if err != nil {
		log.Fatalf("Error posting: %s", err)
	}
	log.Printf("Posted: %s", response.Uri)
}
