package sposter

import (
	"context"
	"fmt"
	"time"

	"github.com/apsvieira/bsky-sposter/src/atproto"
	"github.com/apsvieira/bsky-sposter/src/atproto/richtext"
	"github.com/bluesky-social/indigo/api/bsky"
)

const BLUESKY_SERVICE = "https://bsky.social"

func NewPost(ctx context.Context, client *atproto.Client, text string) (*bsky.FeedPost, error) {
	rt := richtext.NewRichText(text)
	if err := rt.DetectFacets(ctx, client); err != nil {
		return nil, fmt.Errorf("NewPost: %w", err)
	}

	post := &bsky.FeedPost{
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		Text:      rt.Text(),
		Facets:    rt.Facets(),
	}
	return post, nil
}
