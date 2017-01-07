package model

type FullTextSearchResponse struct {
	Users  []*PublicUser `json:"users"`
	Tweets []*Tweet      `json:"tweets"`
}
