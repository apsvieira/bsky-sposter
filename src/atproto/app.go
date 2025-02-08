package atproto

import (
	"context"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/lex/util"
	"github.com/bluesky-social/indigo/xrpc"
)

// AppNS is the namespace for the app.
type AppNS struct {
	client *xrpc.Client
	Bsky   *AppBskyNS
}

func NewAppNS(client *xrpc.Client) *AppNS {
	return &AppNS{client: client, Bsky: NewAppBskyNS(client)}
}

type AppBskyNS struct {
	client *xrpc.Client

	// Actor *AppBskyActorNS
	// Embed *AppBskyEmbedNS
	Feed *AppBskyFeedNS
	// Graph *AppBskyGraphNS
	// Labeler *AppBskyLabelerNS
	// Notification *AppBskyNotificationNS
	// Richtext *AppBskyRichtextNS
	// Unspecced *AppBskyUnspeccedNS
	// Video *AppBskyVideoNS
}

func NewAppBskyNS(client *xrpc.Client) *AppBskyNS {
	return &AppBskyNS{client: client, Feed: NewAppBskyFeedNS(client)}
}

type AppBskyFeedNS struct {
	client *xrpc.Client

	// Generator *GeneratorRecord
	// Like *LikeRecord
	Post *PostRecord
	// Postgate *PostgateRecord
	// Repost *RepostRecord
	// Threadgate *ThreadgateRecord
}

func NewAppBskyFeedNS(client *xrpc.Client) *AppBskyFeedNS {
	return &AppBskyFeedNS{client: client, Post: NewPostRecord(client)}
}

type PostRecord struct {
	client *xrpc.Client
}

func NewPostRecord(client *xrpc.Client) *PostRecord {
	return &PostRecord{client: client}
}

func (p *PostRecord) Create(ctx context.Context, data *bsky.FeedPost) (*atproto.RepoCreateRecord_Output, error) {
	input := &atproto.RepoCreateRecord_Input{
		Collection: "app.bsky.feed.post",
		Repo:       p.client.Auth.Did,
		Record:     &util.LexiconTypeDecoder{Val: data},
	}
	return atproto.RepoCreateRecord(ctx, p.client, input)
}
