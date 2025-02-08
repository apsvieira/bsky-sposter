package main

import (
	"context"
	"log"

	atproto "github.com/apsvieira/bsky-sposter/src/atproto"
	"github.com/apsvieira/bsky-sposter/src/richtext"
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
	log.Printf("Detected facets: %s", rt.Facets[0].Features[0].RichtextFacet_Mention.Did)
	log.Printf("Detected facets: %s", rt.Facets[1].Features[0].RichtextFacet_Link.Uri)
	log.Printf("Detected facets: %s", rt.Facets[2].Features[0].RichtextFacet_Link.Uri)
}
