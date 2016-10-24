package database

import (
	"errors"
	"time"
)

type User struct {
	Id        int64
	Name      string
	Username  string
	Email     string
	CreatedAt time.Time
}

var users = []User{
	{
		Id:        1,
		Name:      "george",
		Username:  "fisher",
		Email:     "corpsegrinder@cannibalcorpse.com",
		CreatedAt: time.Unix(0, 0),
	},
}

func GetUsers() ([]User, error) {
	return users, nil
}

func GetUser(user_id int64) (User, error) {
	user := getUserWithId(user_id)

	/* Emulate DB query fail? */
	if (User{}) == user {
		return User{}, errors.New("User with given ID was not found.")
	}

	return user, nil
}

func InsertUser(user User) (User, error) {
	if userAlreadyExists(user) {
		return User{}, errors.New("User with given username already exists.")
	}

	userId := insertUserToDatabase(user)
	user.Id = userId

	return user, nil
}

/* Functions that mock database queries */

func getUserWithId(userId int64) User {
	var foundUser User

	for _, user := range users {
		if user.Id == userId {
			foundUser = user
		}
	}

	return foundUser
}

func insertUserToDatabase(user User) int64 {
	userId := len(users)
	user.Id = int64(userId)

	users = append(users, user)

	return int64(userId)
}

func userAlreadyExists(userToCheck User) bool {
	for _, user := range users {
		if user.Username == userToCheck.Username {
			return true
		}
	}
	return false
}
