package mock

import (
	"context"
	"fmt"

	"github.com/apsvieira/bsky-sposter/src/atproto/interfaces"
	"github.com/bluesky-social/indigo/api/atproto"
)

func (c *Client) Com() interfaces.ComNS {
	return c.com
}

func (c *Client) App() interfaces.AppNS {

	return c.app
}

type ComNS struct {
	atproto interfaces.ComAtprotoNS
}

func (c *ComNS) Atproto() interfaces.ComAtprotoNS {
	return c.atproto
}

type ComAtprotoNS struct {
	identity interfaces.ComAtprotoIdentityNS
	server   interfaces.ComAtprotoServerNS
}

func (c *ComAtprotoNS) Identity() interfaces.ComAtprotoIdentityNS {
	return c.identity
}

func (c *ComAtprotoNS) Server() interfaces.ComAtprotoServerNS {

	return c.server
}

// ComAtprotoIdentityNS is a mock implementation of ComAtprotoIdentityNS.
type ComAtprotoIdentityNS struct {
}

var _ interfaces.ComAtprotoIdentityNS = (*ComAtprotoIdentityNS)(nil)

func (c *ComAtprotoIdentityNS) ResolveHandle(ctx context.Context, handle string) (*atproto.IdentityResolveHandle_Output, error) {
	return &atproto.IdentityResolveHandle_Output{
		Did: fmt.Sprintf("did:fake:%s", handle),
	}, nil
}
