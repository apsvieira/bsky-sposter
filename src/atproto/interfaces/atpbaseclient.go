package interfaces

type AtpBaseClient interface {
	Com() ComNS
	App() AppNS
	// Chat() ChatNS
	// Feed() FeedNS
}
