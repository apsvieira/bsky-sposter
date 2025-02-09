package mock

import (
	"context"

	"github.com/apsvieira/bsky-sposter/src/atproto/interfaces"
)

type Client struct {
	com interfaces.ComNS
	app interfaces.AppNS
}

func NewClient(ctx context.Context, service string) (*Client, error) {
	client := &Client{
		com: &ComNS{
			atproto: &ComAtprotoNS{
				identity: &ComAtprotoIdentityNS{},
				server:   nil,
			},
		},
		app: nil,
	}
	return client, nil
}
