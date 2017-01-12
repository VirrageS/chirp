package service

import (
	"github.com/VirrageS/chirp/backend/model"
)

// Helper struct to implement sorting of Tweet slices by creation date
type byCreationDateDesc []*model.Tweet

func (s byCreationDateDesc) Len() int {
	return len(s)
}
func (s byCreationDateDesc) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byCreationDateDesc) Less(i, j int) bool {
	return s[i].CreatedAt.After(s[j].CreatedAt)
}
