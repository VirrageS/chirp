package converters

import (
	"database/sql"
	"reflect"
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
		{ // no AvatarURL
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
		{ // no AvatarUrl
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

func TestConvertAPIToDatabase(t *testing.T) {
	// subject
	var converter = NewUserConverter()
	now := time.Now()

	testCases := []struct {
		APIUser *APIModel.NewUserForm
		DBUser  *databaseModel.User
	}{
		{ // positive scenario
			APIUser: &APIModel.NewUserForm{
				Username: "username",
				Password: "password",
				Email:    "user@email.com",
				Name:     "name",
			},
			DBUser: &databaseModel.User{
				ID:            0,
				TwitterToken:  sql.NullString{String: "", Valid: false},
				FacebookToken: sql.NullString{String: "", Valid: false},
				GoogleToken:   sql.NullString{String: "", Valid: false},
				Username:      "username",
				Password:      "password",
				Email:         "user@email.com",
				CreatedAt:     now,
				LastLogin:     now,
				Active:        true,
				Name:          "name",
				AvatarUrl:     sql.NullString{String: "url", Valid: true},
			},
		},
	}

	for _, testCase := range testCases {
		APIUser := testCase.APIUser
		actualDBUser := converter.ConvertAPIToDatabase(APIUser)

		nullString := sql.NullString{String: "", Valid: false}
		// TODO: fix comparison when converter is fixed
		if actualDBUser.ID != 0 ||
			actualDBUser.TwitterToken != nullString ||
			actualDBUser.FacebookToken != nullString ||
			actualDBUser.GoogleToken != nullString ||
			actualDBUser.Username != APIUser.Username ||
			actualDBUser.Password != APIUser.Password ||
			actualDBUser.Email != APIUser.Email ||
			actualDBUser.Name != APIUser.Name ||
			actualDBUser.Active != true ||
			actualDBUser.AvatarUrl != nullString {
			t.Errorf("Got: %v, but expected: %v", actualDBUser, testCase.DBUser)
		}
	}
}

func TestConvertArrayOfDatabaseUser(t *testing.T) {
	// subject
	var converter = NewUserConverter()

	testCases := []struct {
		DBUsers  []*databaseModel.User
		APIUsers []*APIModel.User
	}{
		{ // positive scenario
			DBUsers: []*databaseModel.User{
				{
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
			},
			APIUsers: []*APIModel.User{
				{
					ID:        1,
					Username:  "username",
					Name:      "name",
					AvatarUrl: "url",
					Following: false,
				},
			},
		},
		{ // nil case
			DBUsers:  nil,
			APIUsers: make([]*APIModel.User, 0),
		},
	}

	for _, testCase := range testCases {
		actualAPITweetSlice := converter.ConvertArrayDatabaseToAPI(testCase.DBUsers)
		expectedAPITweetSlice := testCase.APIUsers

		if !reflect.DeepEqual(actualAPITweetSlice, expectedAPITweetSlice) {
			t.Errorf("Got: %v, but expected: %v", actualAPITweetSlice, expectedAPITweetSlice)
		}
	}

}

func TestConvertArrayOfDatabasePublicUser(t *testing.T) {
	// subject
	var converter = NewUserConverter()

	testCases := []struct {
		DBUsers  []*databaseModel.PublicUser
		APIUsers []*APIModel.User
	}{
		{ // positive scenario
			DBUsers: []*databaseModel.PublicUser{
				{
					ID:        1,
					Username:  "username",
					Name:      "name",
					AvatarUrl: sql.NullString{String: "url", Valid: true},
				},
			},
			APIUsers: []*APIModel.User{
				{
					ID:        1,
					Username:  "username",
					Name:      "name",
					AvatarUrl: "url",
					Following: false,
				},
			},
		},
		{ // nil case
			DBUsers:  nil,
			APIUsers: make([]*APIModel.User, 0),
		},
	}

	for _, testCase := range testCases {
		actualAPITweetSlice := converter.ConvertArrayDatabasePublicUserToAPI(testCase.DBUsers)
		expectedAPITweetSlice := testCase.APIUsers

		if !reflect.DeepEqual(actualAPITweetSlice, expectedAPITweetSlice) {
			t.Errorf("Got: %v, but expected: %v", actualAPITweetSlice, expectedAPITweetSlice)
		}
	}
}
