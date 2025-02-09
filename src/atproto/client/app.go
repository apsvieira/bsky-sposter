package client

import (
	"context"

	"github.com/apsvieira/bsky-sposter/src/atproto/interfaces"
	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/xrpc"
)

// AppNS is the namespace for the app.
type AppNS struct {
	client *xrpc.Client
	bsky   interfaces.AppBskyNS
}

var _ interfaces.AppNS = (*AppNS)(nil)

func (a *AppNS) Bsky() interfaces.AppBskyNS {
	return a.bsky
}

var _ interfaces.AppNS = (*AppNS)(nil)

func NewAppNS(client *xrpc.Client) *AppNS {
	return &AppNS{client: client, bsky: NewAppBskyNS(client)}
}

type AppBskyNS struct {
	client *xrpc.Client

	// Actor *AppBskyActorNS
	// Embed *AppBskyEmbedNS
	feed interfaces.AppBskyFeedNS
	// Graph *AppBskyGraphNS
	// Labeler *AppBskyLabelerNS
	// Notification *AppBskyNotificationNS
	// Richtext *AppBskyRichtextNS
	// Unspecced *AppBskyUnspeccedNS
	// Video *AppBskyVideoNS
}

func (a *AppBskyNS) Feed() interfaces.AppBskyFeedNS {
	return a.feed
}

func NewAppBskyNS(client *xrpc.Client) interfaces.AppBskyNS {
	return &AppBskyNS{client: client, feed: NewAppBskyFeedNS(client)}
}

type AppBskyFeedNS struct {
	client *xrpc.Client

	// Generator *GeneratorRecord
	// Like *LikeRecord
	post interfaces.PostRecord
	// Postgate *PostgateRecord
	// Repost *RepostRecord
	// Threadgate *ThreadgateRecord
}

func (a *AppBskyFeedNS) Post() interfaces.PostRecord {
	return a.post
}

func NewAppBskyFeedNS(client *xrpc.Client) interfaces.AppBskyFeedNS {
	return &AppBskyFeedNS{client: client, post: NewPostRecord(client)}
}

type PostRecord struct {
	client *xrpc.Client
}

func NewPostRecord(client *xrpc.Client) interfaces.PostRecord {
	return &PostRecord{client: client}
}

func (p *PostRecord) Create(ctx context.Context, data *atproto.RepoCreateRecord_Input) (*atproto.RepoCreateRecord_Output, error) {
	return atproto.RepoCreateRecord(ctx, p.client, data)
}
