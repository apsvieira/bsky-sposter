package main

import (
	"context"
	"log"

	bsky "github.com/apsvieira/bsky-sposter/src/bsky"
)

func main() {
	ctx := context.Background()
	creds, err := bsky.GetCredentials()
	if err != nil {
		log.Fatalf("Error getting credentials: %s", err)
	}

	client := bsky.NewClient("https://bsky.social", creds)
	if err := client.Authenticate(ctx); err != nil {
		log.Fatalf("Error authenticating: %s", err)
	}
	log.Printf("Authenticated as %s", creds.Handle)
	log.Printf("Client: %#v", client.XrpcClient)
	log.Printf("Auth: %#v", client.XrpcClient.Auth)
}
