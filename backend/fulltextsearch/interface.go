package fulltextsearch

type Searcher interface {
	TweetSearcher
	UserSearcher
}

type TweetSearcher interface {
	GetTweetsIDs(querystring string) ([]int64, error)
}

type UserSearcher interface {
	GetUsersIDs(querystring string) ([]int64, error)
}
