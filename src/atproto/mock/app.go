package mock

import "github.com/apsvieira/bsky-sposter/src/atproto/interfaces"

type AppNS struct {
	bsky interfaces.AppBskyNS
}

func (a *AppNS) Bsky() interfaces.AppBskyNS {
	return a.bsky
}
