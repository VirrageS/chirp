package converters

import (
	"database/sql"
	"testing"
	"time"

	APIModel "github.com/VirrageS/chirp/backend/api/model"
	databaseModel "github.com/VirrageS/chirp/backend/database/model"
)

func TestConvertDatabaseToAPI(t *testing.T) {
	// subject
	var converter = NewUserConverter()

	testCases := []struct {
		DBUser  *databaseModel.User
		APIUser *APIModel.User
	}{
		{ // positive scenario
			DBUser: &databaseModel.User{
				ID:            1,
				TwitterToken:  sql.NullString{String: "", Valid: false},
				FacebookToken: sql.NullString{String: "", Valid: false},
				GoogleToken:   sql.NullString{String: "", Valid: false},
				Username:      "username",
				Password:      "password",
				Email:         "user@email.com",
				CreatedAt:     time.Now(),
				LastLogin:     time.Now(),
				Active:        true,
				Name:          "name",
				AvatarUrl:     sql.NullString{String: "url", Valid: true},
			},
			APIUser: &APIModel.User{
				ID:        1,
				Username:  "username",
				Name:      "name",
				AvatarUrl: "url",
				Following: false,
			},
		},
		{ // invalid AvatarURL
			DBUser: &databaseModel.User{
				ID:            1,
				TwitterToken:  sql.NullString{String: "", Valid: false},
				FacebookToken: sql.NullString{String: "", Valid: false},
				GoogleToken:   sql.NullString{String: "", Valid: false},
				Username:      "username",
				Password:      "password",
				Email:         "user@email.com",
				CreatedAt:     time.Now(),
				LastLogin:     time.Now(),
				Active:        true,
				Name:          "name",
				AvatarUrl:     sql.NullString{String: "url", Valid: false},
			},
			APIUser: &APIModel.User{
				ID:        1,
				Username:  "username",
				Name:      "name",
				AvatarUrl: "",
				Following: false,
			},
		},
	}

	for _, testCase := range testCases {
		actualApiUser := converter.ConvertDatabaseToAPI(testCase.DBUser)

		if *actualApiUser != *testCase.APIUser {
			t.Errorf("Got: %v, but expected: %v", actualApiUser, testCase.APIUser)
		}
	}
}

func TestConvertDatabasePublicUserToAPI(t *testing.T) {
	// subject
	var converter = NewUserConverter()

	testCases := []struct {
		DBUser  *databaseModel.PublicUser
		APIUser *APIModel.User
	}{
		{ // positive scenario
			DBUser: &databaseModel.PublicUser{
				ID:        1,
				Username:  "username",
				Name:      "name",
				AvatarUrl: sql.NullString{String: "url", Valid: true},
			},
			APIUser: &APIModel.User{
				ID:        1,
				Username:  "username",
				Name:      "name",
				AvatarUrl: "url",
				Following: false,
			},
		},
		{ // invalid AvatarUrl
			DBUser: &databaseModel.PublicUser{
				ID:        1,
				Username:  "username",
				Name:      "name",
				AvatarUrl: sql.NullString{String: "url", Valid: false},
			},
			APIUser: &APIModel.User{
				ID:        1,
				Username:  "username",
				Name:      "name",
				AvatarUrl: "",
				Following: false,
			},
		},
	}

	for _, testCase := range testCases {
		actualAPIUser := converter.ConvertDatabasePublicUserToAPI(testCase.DBUser)

		if *actualAPIUser != *testCase.APIUser {
			t.Errorf("Got: %v, but expected: %v", actualAPIUser, testCase.APIUser)
		}
	}
}
