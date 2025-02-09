package main

import (
	"context"
	"log"

	sposter "github.com/apsvieira/bsky-sposter/src"
	"github.com/apsvieira/bsky-sposter/src/atproto/client"
)

func main() {
	ctx := context.Background()
	creds, err := client.GetCredentials()
	if err != nil {
		log.Fatalf("Error getting credentials: %s", err)
	}

	client, err := client.NewClient(ctx, "https://bsky.social", creds)
	if err != nil {
		log.Fatalf("Error authenticating: %s", err)
	}
	log.Printf("Authenticated as %s", creds.Handle)

	post, err := sposter.NewPost(ctx, client, "Hello, @apsv.bsky.social! #bsky https://bsky.social! google.com ")
	if err != nil {
		log.Fatalf("Error creating post: %s", err)
	}

	bskyPost, err := client.CreatePost(ctx, post)
	if err != nil {
		log.Fatalf("Error posting: %s", err)
	}
	log.Printf("Posted: %s", bskyPost.Uri)
}
