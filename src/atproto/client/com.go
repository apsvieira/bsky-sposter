package client

import (
	"context"
	"fmt"

	"github.com/apsvieira/bsky-sposter/src/atproto/interfaces"
	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/xrpc"
)

type ComNS struct {
	client  *xrpc.Client
	atproto interfaces.ComAtprotoNS
}

var _ interfaces.ComNS = (*ComNS)(nil)

func (c *ComNS) Atproto() interfaces.ComAtprotoNS {
	return c.atproto
}

func NewComNS(client *xrpc.Client) *ComNS {
	return &ComNS{client: client, atproto: NewComAtprotoNS(client)}
}

type ComAtprotoNS struct {
	client *xrpc.Client
	// Admin ComAtprotoAdminNS
	identity interfaces.ComAtprotoIdentityNS
	// Label *ComAtprotoLabelNS
	// Lexicon *ComAtprotoLexiconNS
	// Moderation *ComAtprotoModerationNS
	// Repo *ComAtprotoRepoNS
	server interfaces.ComAtprotoServerNS
	// Sync *ComAtprotoSyncNS
	// Temp *ComAtprotoTempNS
}

func (c ComAtprotoNS) Identity() interfaces.ComAtprotoIdentityNS {
	return c.identity
}

func (c ComAtprotoNS) Server() interfaces.ComAtprotoServerNS {
	return c.server
}

func NewComAtprotoNS(client *xrpc.Client) interfaces.ComAtprotoNS {
	return &ComAtprotoNS{
		client:   client,
		identity: NewComAtprotoIdentityNS(client),
		server:   NewComAtprotoServerNS(client),
	}
}

type ComAtprotoIdentityNS struct {
	client *xrpc.Client
}

var _ interfaces.ComAtprotoIdentityNS = (*ComAtprotoIdentityNS)(nil)

func NewComAtprotoIdentityNS(client *xrpc.Client) interfaces.ComAtprotoIdentityNS {
	return &ComAtprotoIdentityNS{client: client}
}

func (c ComAtprotoIdentityNS) ResolveHandle(ctx context.Context, handle string) (*atproto.IdentityResolveHandle_Output, error) {
	return atproto.IdentityResolveHandle(ctx, c.client, handle)
}

type ComAtprotoServerNS struct {
	client *xrpc.Client
}

var _ interfaces.ComAtprotoServerNS = (*ComAtprotoServerNS)(nil)

func NewComAtprotoServerNS(client *xrpc.Client) interfaces.ComAtprotoServerNS {
	return &ComAtprotoServerNS{client: client}
}

// Authenticate creates a new session and authenticates with the handle and appkey.
func (c *ComAtprotoServerNS) CreateSession(ctx context.Context, data *atproto.ServerCreateSession_Input) error {
	session, err := atproto.ServerCreateSession(ctx, c.client, data)
	if err != nil {
		return fmt.Errorf("Authenticate: %w", err)
	}
	if !*session.Active {
		return fmt.Errorf("Authenticate: user not active: %v", *session.Status)
	}

	c.client.Auth = &xrpc.AuthInfo{
		Did:        session.Did,
		Handle:     session.Handle,
		AccessJwt:  session.AccessJwt,
		RefreshJwt: session.RefreshJwt,
	}
	return nil
}

// RefreshSession refreshes the current session, creating a new access token.
func (c *ComAtprotoServerNS) RefreshSession(ctx context.Context) error {
	if c.client.Auth == nil || c.client.Auth.RefreshJwt == "" {
		return fmt.Errorf("RefreshSession: no session to refresh")
	}

	session, err := atproto.ServerRefreshSession(ctx, c.client)
	if err != nil {
		return fmt.Errorf("RefreshSession: %w", err)
	}
	if !*session.Active {
		return fmt.Errorf("RefreshSession: user not active: %v", *session.Status)
	}

	c.client.Auth = &xrpc.AuthInfo{
		Did:        session.Did,
		Handle:     session.Handle,
		AccessJwt:  session.AccessJwt,
		RefreshJwt: session.RefreshJwt,
	}
	return nil
}
