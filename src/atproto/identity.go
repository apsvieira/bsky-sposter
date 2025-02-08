package atproto

import (
	"context"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/xrpc"
)

type ComAtprotoIdentityNS struct {
	client *xrpc.Client
}

func NewComAtprotoIdentityNS(client *xrpc.Client) *ComAtprotoIdentityNS {
	return &ComAtprotoIdentityNS{client: client}
}

func (c *ComAtprotoIdentityNS) ResolveHandle(ctx context.Context, handle string) (*atproto.IdentityResolveHandle_Output, error) {
	return atproto.IdentityResolveHandle(ctx, c.client, handle)
}
