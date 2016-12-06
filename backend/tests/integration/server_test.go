package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"gopkg.in/gin-gonic/gin.v1"

	"github.com/VirrageS/chirp/backend/database"
	"github.com/VirrageS/chirp/backend/model"
	"github.com/VirrageS/chirp/backend/server"
	"github.com/stretchr/testify/assert"
)

const baseURL = "http://localhost:8080"

var s *gin.Engine

func TestMain(m *testing.M) {
	db := database.NewConnection("5555")
	db.Exec("TRUNCATE users, tweets CASCADE;") // Ugly, but lets keep it for convenience for now
	gin.SetMode(gin.TestMode)
	s = server.New(db)
	os.Exit(m.Run())
}

func TestCreateNewUser(t *testing.T) {
	newUser := &model.NewUserForm{
		Username: "user",
		Password: "password",
		Email:    "email@email.com",
		Name:     "name",
	}
	data, _ := json.Marshal(newUser)
	buf := bytes.NewBuffer(data)

	req, _ := http.NewRequest("POST", baseURL+"/signup", buf)
	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	var actualUser model.PublicUser
	err := json.Unmarshal(w.Body.Bytes(), &actualUser)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, actualUser.Username, newUser.Username)
	assert.Equal(t, actualUser.Name, newUser.Name)
	assert.Equal(t, actualUser.AvatarUrl, "")
	assert.Equal(t, actualUser.Following, false)
}

func TestCreateUserWithUsernameThatAlreadyExists(t *testing.T) {
	user := &model.NewUserForm{
		Username: "user10",
		Password: "password10",
		Email:    "email10@email.com",
		Name:     "name10",
	}
	_, err := createUser(user, s, baseURL)
	assert.Nil(t, err)

	newUser := &model.NewUserForm{
		Username: "user10",
		Password: "password11",
		Email:    "email11@email.com",
		Name:     "name11",
	}

	// try to create the same user again
	data, _ := json.Marshal(newUser)
	buf := bytes.NewBuffer(data)

	req, _ := http.NewRequest("POST", baseURL+"/signup", buf)
	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	// TODO: add ID check once database is cleaned after every run
	assert.Equal(t, http.StatusConflict, w.Code)
	assert.NotEmpty(t, w.Body)
}

func TestCreateUserWithEmailThatAlreadyExists(t *testing.T) {
	user := &model.NewUserForm{
		Username: "user20",
		Password: "password20",
		Email:    "email20@email.com",
		Name:     "name20",
	}
	_, err := createUser(user, s, baseURL)
	assert.Nil(t, err)

	newUser := &model.NewUserForm{
		Username: "user21",
		Password: "password21",
		Email:    "email20@email.com",
		Name:     "name21",
	}

	// try to create the same user again
	data, _ := json.Marshal(newUser)
	buf := bytes.NewBuffer(data)

	req, _ := http.NewRequest("POST", baseURL+"/signup", buf)
	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	// TODO: add ID check once database is cleaned after every run
	assert.Equal(t, http.StatusConflict, w.Code)
	assert.NotEmpty(t, w.Body)
}

func TestLoginUser(t *testing.T) {
	// TODO: maybe prepare test user at the beginning of tests?
	newUser := &model.NewUserForm{
		Username: "user1",
		Password: "password1",
		Email:    "email1@email.com",
		Name:     "name1",
	}
	userID, err := createUser(newUser, s, baseURL)
	assert.Nil(t, err)

	loginData := &model.LoginForm{
		Email:    "email1@email.com",
		Password: "password1",
	}
	data, _ := json.Marshal(loginData)
	buf := bytes.NewBuffer(data)

	req, _ := http.NewRequest("POST", baseURL+"/login", buf)
	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	var loginResponse model.LoginResponse
	err = json.Unmarshal(w.Body.Bytes(), &loginResponse)
	assert.Nil(t, err)

	actualUser := loginResponse.User

	expectedUser := &model.PublicUser{
		ID:        userID,
		Username:  newUser.Username,
		Name:      newUser.Name,
		AvatarUrl: "",
		Following: false,
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedUser, actualUser)
}

func TestLoginUserWithInvalidPassword(t *testing.T) {
	// TODO: maybe prepare test user at the beginning of tests?
	newUser := &model.NewUserForm{
		Username: "user40",
		Password: "password40",
		Email:    "email40@email.com",
		Name:     "name40",
	}
	_, err := createUser(newUser, s, baseURL)
	assert.Nil(t, err)

	loginData := &model.LoginForm{
		Email:    "email40@email.com",
		Password: "invalidpassword",
	}
	data, _ := json.Marshal(loginData)
	buf := bytes.NewBuffer(data)

	req, _ := http.NewRequest("POST", baseURL+"/login", buf)
	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.NotEmpty(t, w.Body)
}

func TestLoginUserWithInvalidEmail(t *testing.T) {
	// TODO: maybe prepare test user at the beginning of tests?
	newUser := &model.NewUserForm{
		Username: "user50",
		Password: "password50",
		Email:    "email50@email.com",
		Name:     "name50",
	}
	_, err := createUser(newUser, s, baseURL)
	assert.Nil(t, err)

	loginData := &model.LoginForm{
		Email:    "invalid@email.com",
		Password: "password50",
	}
	data, _ := json.Marshal(loginData)
	buf := bytes.NewBuffer(data)

	req, _ := http.NewRequest("POST", baseURL+"/login", buf)
	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.NotEmpty(t, w.Body)
}

func TestCreateTweetResponse(t *testing.T) {
	// TODO: maybe prepare test user at the beginning of tests?
	newUser := &model.NewUserForm{
		Username: "user2",
		Password: "password2",
		Email:    "email2@email.com",
		Name:     "name2",
	}
	userID, err := createUser(newUser, s, baseURL)
	assert.Nil(t, err)

	loginData := &model.LoginForm{
		Email:    "email2@email.com",
		Password: "password2",
	}
	authToken, err := loginUser(loginData, s, baseURL)
	assert.Nil(t, err)

	newTweet := &model.NewTweet{
		Content: "new tweet",
	}
	data, _ := json.Marshal(newTweet)
	buf := bytes.NewBuffer(data)

	req, _ := http.NewRequest("POST", baseURL+"/tweets", buf)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+authToken)

	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	var actualTweet model.Tweet
	err = json.Unmarshal(w.Body.Bytes(), &actualTweet)
	assert.Nil(t, err)

	expectedUser := &model.PublicUser{
		ID:        userID,
		Username:  newUser.Username,
		Name:      newUser.Name,
		AvatarUrl: "",
		Following: false,
	}

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, int64(0), actualTweet.Likes)
	assert.Equal(t, int64(0), actualTweet.Retweets)
	assert.Equal(t, "new tweet", actualTweet.Content)
	assert.Equal(t, false, actualTweet.Liked)
	assert.Equal(t, false, actualTweet.Retweeted)
	assert.Equal(t, expectedUser, actualTweet.Author)
}

func TestGetTweetAfterCreatingTweet(t *testing.T) {
	// TODO: maybe prepare test user at the beginning of tests?
	newUser := &model.NewUserForm{
		Username: "user3",
		Password: "password3",
		Email:    "email3@email.com",
		Name:     "name3",
	}
	_, err := createUser(newUser, s, baseURL)
	assert.Nil(t, err)

	loginData := &model.LoginForm{
		Email:    "email3@email.com",
		Password: "password3",
	}
	authToken, err := loginUser(loginData, s, baseURL)
	assert.Nil(t, err)

	newTweet := &model.NewTweet{
		Content: "new tweet",
	}
	data, _ := json.Marshal(newTweet)
	buf := bytes.NewBuffer(data)

	req, _ := http.NewRequest("POST", baseURL+"/tweets", buf)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+authToken)

	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	var createdTweet model.Tweet
	err = json.Unmarshal(w.Body.Bytes(), &createdTweet)
	assert.Nil(t, err)

	tweetID := createdTweet.ID

	reqGET, _ := http.NewRequest("GET", baseURL+"/tweets/"+strconv.FormatInt(int64(tweetID), 10), buf)
	reqGET.Header.Add("Authorization", "Bearer "+authToken)

	w2 := httptest.NewRecorder()
	s.ServeHTTP(w2, reqGET)

	var actualTweet model.Tweet
	err = json.Unmarshal(w2.Body.Bytes(), &actualTweet)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Equal(t, createdTweet, actualTweet)
}

func TestGetTweetsAfterCreatingTweets(t *testing.T) {
	// TODO: maybe prepare test user at the beginning of tests?
	newUser := &model.NewUserForm{
		Username: "user33",
		Password: "password33",
		Email:    "email33@email.com",
		Name:     "name33",
	}
	_, err := createUser(newUser, s, baseURL)
	assert.Nil(t, err)

	loginData := &model.LoginForm{
		Email:    "email33@email.com",
		Password: "password33",
	}
	authToken, err := loginUser(loginData, s, baseURL)
	assert.Nil(t, err)

	// create frist tweet
	newTweet1 := &model.NewTweet{
		Content: "new tweet1",
	}
	data, _ := json.Marshal(newTweet1)
	buf := bytes.NewBuffer(data)

	req, _ := http.NewRequest("POST", baseURL+"/tweets", buf)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+authToken)

	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	var createdTweet1 model.Tweet
	err = json.Unmarshal(w.Body.Bytes(), &createdTweet1)
	assert.Nil(t, err)

	// create another tweet
	newTweet2 := &model.NewTweet{
		Content: "new tweet1",
	}
	data, _ = json.Marshal(newTweet2)
	buf = bytes.NewBuffer(data)

	req2, _ := http.NewRequest("POST", baseURL+"/tweets", buf)
	req2.Header.Add("Content-Type", "application/json")
	req2.Header.Add("Authorization", "Bearer "+authToken)

	w2 := httptest.NewRecorder()
	s.ServeHTTP(w2, req2)

	var createdTweet2 model.Tweet
	err = json.Unmarshal(w2.Body.Bytes(), &createdTweet2)
	assert.Nil(t, err)

	reqGET, _ := http.NewRequest("GET", baseURL+"/tweets", buf)
	reqGET.Header.Add("Authorization", "Bearer "+authToken)

	w3 := httptest.NewRecorder()
	s.ServeHTTP(w3, reqGET)

	var actualTweets []model.Tweet
	err = json.Unmarshal(w3.Body.Bytes(), &actualTweets)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, w3.Code)
	assert.Contains(t, actualTweets, createdTweet1)
	assert.Contains(t, actualTweets, createdTweet2)
}

func TestDeleteTweetAfterCreatingTweet(t *testing.T) {
	// TODO: maybe prepare test user at the beginning of tests?
	newUser := &model.NewUserForm{
		Username: "user35",
		Password: "password35",
		Email:    "email35@email.com",
		Name:     "name35",
	}
	_, err := createUser(newUser, s, baseURL)
	assert.Nil(t, err)

	loginData := &model.LoginForm{
		Email:    "email35@email.com",
		Password: "password35",
	}
	authToken, err := loginUser(loginData, s, baseURL)
	assert.Nil(t, err)

	newTweet := &model.NewTweet{
		Content: "new tweet",
	}
	data, _ := json.Marshal(newTweet)
	buf := bytes.NewBuffer(data)

	req, _ := http.NewRequest("POST", baseURL+"/tweets", buf)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+authToken)

	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	var createdTweet model.Tweet
	err = json.Unmarshal(w.Body.Bytes(), &createdTweet)
	assert.Nil(t, err)

	tweetID := createdTweet.ID

	reqDELETE, _ := http.NewRequest("DELETE", baseURL+"/tweets/"+strconv.FormatInt(int64(tweetID), 10), buf)
	reqDELETE.Header.Add("Authorization", "Bearer "+authToken)

	w2 := httptest.NewRecorder()
	s.ServeHTTP(w2, reqDELETE)

	assert.Equal(t, http.StatusNoContent, w2.Code)

	reqGET, _ := http.NewRequest("GET", baseURL+"/tweets/"+strconv.FormatInt(int64(tweetID), 10), buf)
	reqGET.Header.Add("Authorization", "Bearer "+authToken)

	w3 := httptest.NewRecorder()
	s.ServeHTTP(w3, reqGET)

	assert.Equal(t, http.StatusNotFound, w3.Code)
}
