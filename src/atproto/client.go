package atproto

import (
	"context"
	"fmt"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/xrpc"
)

type Client struct {
	client *xrpc.Client
	creds  *Credentials
	Com    *ComNS
}

func NewClient(ctx context.Context, service string, creds *Credentials) (*Client, error) {
	xrpcClient := &xrpc.Client{
		Host: service,
	}
	client := &Client{client: xrpcClient, creds: creds, Com: NewComNS(xrpcClient)}
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

	session, err := atproto.ServerCreateSession(ctx, c.client, creds)
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
func (c *Client) RefreshSession(ctx context.Context) error {
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

type ComNS struct {
	client  *xrpc.Client
	Atproto *ComAtprotoNS
}

func NewComNS(client *xrpc.Client) *ComNS {
	return &ComNS{client: client, Atproto: NewComAtprotoNS(client)}
}

type ComAtprotoNS struct {
	client *xrpc.Client
	// Admin ComAtprotoAdminNS
	Identity *ComAtprotoIdentityNS
	// Label *ComAtprotoLabelNS
	// Lexicon *ComAtprotoLexiconNS
	// Moderation *ComAtprotoModerationNS
	// Repo *ComAtprotoRepoNS
	// Server *ComAtprotoServerNS
	// Sync *ComAtprotoSyncNS
	// Temp *ComAtprotoTempNS
}

func NewComAtprotoNS(client *xrpc.Client) *ComAtprotoNS {
	return &ComAtprotoNS{
		client:   client,
		Identity: NewComAtprotoIdentityNS(client),
	}
}
