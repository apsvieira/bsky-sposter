package bsky

import (
	"context"
	"fmt"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/xrpc"
)

type Client struct {
	XrpcClient *xrpc.Client
	creds      *Credentials
}

func NewClient(service string, creds *Credentials) *Client {
	xrpcClient := &xrpc.Client{
		Host: service,
	}
	return &Client{XrpcClient: xrpcClient, creds: creds}
}

// Authenticate creates a new session and authenticates with the handle and appkey.
func (c *Client) Authenticate(ctx context.Context) error {
	creds := &atproto.ServerCreateSession_Input{
		Identifier: c.creds.Handle,
		Password:   c.creds.AppKey,
	}

	session, err := atproto.ServerCreateSession(ctx, c.XrpcClient, creds)
	if err != nil {
		return fmt.Errorf("Authenticate: %w", err)
	}
	if !*session.Active {
		return fmt.Errorf("Authenticate: session not active: %v", *session.Status)
	}

	c.XrpcClient.Auth = &xrpc.AuthInfo{
		Did:        session.Did,
		Handle:     session.Handle,
		AccessJwt:  session.AccessJwt,
		RefreshJwt: session.RefreshJwt,
	}
	return nil
}

func (c *Client) RefreshSession(ctx context.Context) error {
	if c.XrpcClient.Auth == nil || c.XrpcClient.Auth.RefreshJwt == "" {
		return fmt.Errorf("RefreshSession: no session to refresh")
	}

	session, err := atproto.ServerRefreshSession(ctx, c.XrpcClient)
	if err != nil {
		return fmt.Errorf("RefreshSession: %w", err)
	}
	if !*session.Active {
		return fmt.Errorf("RefreshSession: session not active: %v", *session.Status)
	}

	c.XrpcClient.Auth = &xrpc.AuthInfo{
		Did:        session.Did,
		Handle:     session.Handle,
		AccessJwt:  session.AccessJwt,
		RefreshJwt: session.RefreshJwt,
	}
	return nil
}
