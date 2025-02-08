package atproto

import (
	"context"
	"fmt"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/xrpc"
)

type ComAtprotoServerNS struct {
	client *xrpc.Client
}

func NewComAtprotoServerNS(client *xrpc.Client) *ComAtprotoServerNS {
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
