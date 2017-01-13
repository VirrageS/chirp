package utils

import (
	"github.com/VirrageS/chirp/backend/model"
)

// Helper struct to implement sorting of Tweet slices by creation date
type TweetsByCreationDateDesc []*model.Tweet

func (s TweetsByCreationDateDesc) Len() int {
	return len(s)
}
func (s TweetsByCreationDateDesc) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s TweetsByCreationDateDesc) Less(i, j int) bool {
	return s[i].CreatedAt.After(s[j].CreatedAt)
}
