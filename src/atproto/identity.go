package atproto

import (
	"context"
	"fmt"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/xrpc"
)

type ComAtprotoIdentityNS interface {
	ResolveHandle(ctx context.Context, handle string) (*atproto.IdentityResolveHandle_Output, error)
}

type AtprotoIdentityNS struct {
	client *xrpc.Client
}

var _ ComAtprotoIdentityNS = (*AtprotoIdentityNS)(nil)

func NewComAtprotoIdentityNS(client *xrpc.Client) ComAtprotoIdentityNS {
	return &AtprotoIdentityNS{client: client}
}

func (c *AtprotoIdentityNS) ResolveHandle(ctx context.Context, handle string) (*atproto.IdentityResolveHandle_Output, error) {
	return atproto.IdentityResolveHandle(ctx, c.client, handle)
}

// MockAtprotoIdentityNS is a mock implementation of ComAtprotoIdentityNS.
type MockAtprotoIdentityNS struct {
}

var _ ComAtprotoIdentityNS = (*MockAtprotoIdentityNS)(nil)

func (c *MockAtprotoIdentityNS) ResolveHandle(ctx context.Context, handle string) (*atproto.IdentityResolveHandle_Output, error) {
	return &atproto.IdentityResolveHandle_Output{
		Did: fmt.Sprintf("did:fake:%s", handle),
	}, nil
}
