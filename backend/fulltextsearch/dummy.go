package fulltextsearch

type DummySearch struct{}

func NewDummySearch() Searcher {
	return &DummySearch{}
}

func (d *DummySearch) GetTweetsIDs(querystring string) ([]int64, error) {
	return nil, nil
}

func (d *DummySearch) GetUsersIDs(querystring string) ([]int64, error) {
	return nil, nil
}
