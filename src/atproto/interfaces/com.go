package interfaces

import (
	"context"

	"github.com/bluesky-social/indigo/api/atproto"
)

type ComNS interface {
	Atproto() ComAtprotoNS
}

type ComAtprotoNS interface {
	// Admin ComAtprotoAdminNS
	Identity() ComAtprotoIdentityNS
	// Label *ComAtprotoLabelNS
	// Lexicon *ComAtprotoLexiconNS
	// Moderation *ComAtprotoModerationNS
	// Repo *ComAtprotoRepoNS
	Server() ComAtprotoServerNS
	// Sync *ComAtprotoSyncNS
	// Temp *ComAtprotoTempNS
}

type ComAtprotoIdentityNS interface {
	ResolveHandle(ctx context.Context, handle string) (*atproto.IdentityResolveHandle_Output, error)
}

type ComAtprotoServerNS interface {
	CreateSession(ctx context.Context, data *atproto.ServerCreateSession_Input) error
	RefreshSession(ctx context.Context) error
}
