package client

import (
	"context"
	"fmt"

	"github.com/apsvieira/bsky-sposter/src/atproto/interfaces"
	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/lex/util"
	"github.com/bluesky-social/indigo/xrpc"
)

type Client struct {
	client *xrpc.Client
	creds  *Credentials
	com    interfaces.ComNS
	app    interfaces.AppNS
}

var _ interfaces.AtpBaseClient = (*Client)(nil)

func (c *Client) Com() interfaces.ComNS {
	return c.com
}

func (c *Client) App() interfaces.AppNS {
	return c.app
}
func NewClient(ctx context.Context, service string, creds *Credentials) (*Client, error) {
	xrpcClient := &xrpc.Client{
		Host: service,
	}
	client := &Client{
		client: xrpcClient,
		creds:  creds,
		com:    NewComNS(xrpcClient),
		app:    NewAppNS(xrpcClient),
	}

	if err := client.Authenticate(ctx); err != nil {
		return nil, fmt.Errorf("NewClient: %w", err)
	}
	return client, nil
}

// Authenticate creates a new session and authenticates with the handle and appkey.
func (c *Client) Authenticate(ctx context.Context) error {
	creds := &atproto.ServerCreateSession_Input{
		Identifier: c.creds.Handle,
		Password:   c.creds.AppKey,
	}

	err := c.Com().Atproto().Server().CreateSession(ctx, creds)
	if err != nil {
		return fmt.Errorf("Authenticate: %w", err)
	}
	return nil
}

func (c *Client) CreatePost(ctx context.Context, data *bsky.FeedPost) (*atproto.RepoCreateRecord_Output, error) {
	input := &atproto.RepoCreateRecord_Input{
		Collection: "app.bsky.feed.post",
		Repo:       c.client.Auth.Did,
		Record:     &util.LexiconTypeDecoder{Val: data},
	}
	return atproto.RepoCreateRecord(ctx, c.client, input)
}
