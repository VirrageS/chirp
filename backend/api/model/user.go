package model

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Name      string `json:"name"`
	AvatarUrl string `json:"avatar_url"`
	Following bool   `json:"following"`
}
