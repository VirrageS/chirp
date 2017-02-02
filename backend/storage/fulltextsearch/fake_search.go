package fulltextsearch

type fakeSearch struct{}

// NewFakeSearch creates new instance of fake searcher which imitates searching values.
// Implements all Searcher methods but each of them is doing nothing.
func NewFakeSearch() Searcher {
	return &fakeSearch{}
}

func (d *fakeSearch) GetTweetsIDs(querystring string) ([]int64, error) {
	return nil, nil
}

func (d *fakeSearch) GetUsersIDs(querystring string) ([]int64, error) {
	return nil, nil
}
