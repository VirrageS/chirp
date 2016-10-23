package database

import (
	"errors"
	"time"
)

type User struct {
	Id 		int64
	Name 		string
	Username 	string
	Email	 	string
	CreatedAt 	time.Time
}

// TODO: replace with real databse queries

var users = []User{
	{
		Id: 1,
		Name: "george",
		Username: "fisher",
		Email: "corpsegrinder@cannibalcorpse.com",
		CreatedAt: time.Unix(0, 0),
	},
}

func GetUser(user_id int64) (User, error) {
	user := getUserWithId(user_id)

	/* Emulate DB query fail? */
	if (User{}) == user {
		return User{}, errors.New("User with given ID was not found.")
	}

	return user, nil
}

// TODO: rewrite to create new object instead of editing function parameter?
func InsertUser(user User) (User, error) {
	if userAlreadyExists(user) {
		return User{}, errors.New("User with given username already exists.")
	}

	user_id := insertUserToDatabase(user)
	user.Id = user_id

	return user, nil
}

func getUserWithId(user_id int64) User {
	var found_user User

	for _, user := range users {
		if user.Id == user_id {
			found_user = user
		}
	}

	return found_user
}

func insertUserToDatabase(user User) int64 {
	user_id := len(users)
	user.Id = int64(user_id)

	users = append(users, user)

	return int64(user_id)
}

func userAlreadyExists(user_to_check User) bool {
	for _, user := range users {
		if user.Username == user_to_check.Username {
			return true
		}
	}
	return false
}
