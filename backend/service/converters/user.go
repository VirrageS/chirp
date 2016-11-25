package converters

// TODO: refactor those converters -> maybe they can be more 'generic', because at the moment
// publicUser and User converters are EXACTLY the same

import (
	"time"

	APIModel "github.com/VirrageS/chirp/backend/api/model"
	databaseModel "github.com/VirrageS/chirp/backend/database/model"
)

type UserModelConverter interface {
	ConvertAPIToDatabase(user *APIModel.NewUserForm) *databaseModel.User
	ConvertDatabaseToAPI(user *databaseModel.User) *APIModel.User
	ConvertDatabasePublicUserToAPI(user *databaseModel.PublicUser) *APIModel.User
	ConvertArrayDatabaseToAPI(databaseUsers []*databaseModel.User) []*APIModel.User
	ConvertArrayDatabasePublicUserToAPI(databaseUsers []*databaseModel.PublicUser) []*APIModel.User
}

type UserConverter struct{}

func NewUserConverter() *UserConverter {
	return &UserConverter{}
}

func (converter *UserConverter) ConvertDatabaseToAPI(user *databaseModel.User) *APIModel.User {
	id := user.ID
	username := user.Username
	name := user.Name
	var avatarURL string

	if user.AvatarUrl.Valid {
		avatarURL = user.AvatarUrl.String
	} else {
		avatarURL = ""
	}

	return &APIModel.User{
		ID:        id,
		Username:  username,
		Name:      name,
		AvatarUrl: avatarURL,
		Following: false,
	}
}

//TODO: This should not use time.Now(), but should get an object that stores the time
func (converter *UserConverter) ConvertAPIToDatabase(user *APIModel.NewUserForm) *databaseModel.User {
	username := user.Username
	password := user.Password
	email := user.Email
	name := user.Name
	creationTime := time.Now()

	return &databaseModel.User{
		ID:            0,
		TwitterToken:  toSqlNullString(""),
		FacebookToken: toSqlNullString(""),
		GoogleToken:   toSqlNullString(""),
		Username:      username,
		Password:      password,
		Email:         email,
		CreatedAt:     creationTime,
		LastLogin:     creationTime,
		Active:        true,
		Name:          name,
		AvatarUrl:     toSqlNullString(""),
	}
}

func (converter *UserConverter) ConvertDatabasePublicUserToAPI(user *databaseModel.PublicUser) *APIModel.User {
	id := user.ID
	username := user.Username
	name := user.Name
	var avatarURL string

	if user.AvatarUrl.Valid {
		avatarURL = user.AvatarUrl.String
	} else {
		avatarURL = ""
	}

	return &APIModel.User{
		ID:        id,
		Username:  username,
		Name:      name,
		AvatarUrl: avatarURL,
		Following: false,
	}
}

func (converter *UserConverter) ConvertArrayDatabaseToAPI(databaseUsers []*databaseModel.User) []*APIModel.User {
	convertedUsers := make([]*APIModel.User, 0)

	for _, databaseUser := range databaseUsers {
		APIUser := converter.ConvertDatabaseToAPI(databaseUser)
		convertedUsers = append(convertedUsers, APIUser)
	}

	return convertedUsers
}

func (converter *UserConverter) ConvertArrayDatabasePublicUserToAPI(databaseUsers []*databaseModel.PublicUser) []*APIModel.User {
	convertedUsers := make([]*APIModel.User, 0)

	for _, databaseUser := range databaseUsers {
		APIUser := converter.ConvertDatabasePublicUserToAPI(databaseUser)
		convertedUsers = append(convertedUsers, APIUser)
	}

	return convertedUsers
}
