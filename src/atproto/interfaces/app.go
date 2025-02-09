package interfaces

import (
	"context"

	"github.com/bluesky-social/indigo/api/atproto"
)

type AppNS interface {
	Bsky() AppBskyNS
}

type AppBskyNS interface {
	Feed() AppBskyFeedNS
}

type AppBskyFeedNS interface {
	// Generator *GeneratorRecord
	// Like *LikeRecord
	Post() PostRecord
	// Postgate *PostgateRecord
	// Repost *RepostRecord
	// Threadgate *ThreadgateRecord
}

type PostRecord interface {
	Create(ctx context.Context, data *atproto.RepoCreateRecord_Input) (*atproto.RepoCreateRecord_Output, error)
}
