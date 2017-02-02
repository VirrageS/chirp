package fulltextsearch

// TweetsSearcher is interface which defines all full text search functions which
// are connected with tweets.
type TweetsSearcher interface {
	GetTweetsIDs(querystring string) ([]int64, error)
}

// UsersSearcher is interface which defines all full text search functions which
// are connected with users.
type UsersSearcher interface {
	GetUsersIDs(querystring string) ([]int64, error)
}

// Searcher is interface which defines all full text search functions used in system.
// All implementations of different types of search engines should implement these methods.
type Searcher interface {
	TweetsSearcher
	UsersSearcher
}
