package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/VirrageS/chirp/backend/model"
)

var baseURL string
var s *gin.Engine

var testUser model.User
var otherTestUser model.User

func TestMain(m *testing.M) {
	setup(&testUser, &otherTestUser, &s, baseURL)
	os.Exit(m.Run())
}

func TestCreateNewUser(t *testing.T) {
	newUser := &model.NewUserForm{
		Username: "anotherUser",
		Password: "anotherPassword",
		Email:    "another@email.com",
		Name:     "anotherName",
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
	newUser := &model.NewUserForm{
		Username: testUser.Username,
		Password: "somepassword",
		Email:    "some@email.com",
		Name:     "somename",
	}

	data, _ := json.Marshal(newUser)
	buf := bytes.NewBuffer(data)

	req, _ := http.NewRequest("POST", baseURL+"/signup", buf)
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.NotEmpty(t, w.Body)
}

func TestCreateUserWithEmailThatAlreadyExists(t *testing.T) {
	newUser := &model.NewUserForm{
		Username: "someusername",
		Password: "somepassword",
		Email:    testUser.Email,
		Name:     "somename",
	}

	data, _ := json.Marshal(newUser)
	buf := bytes.NewBuffer(data)

	req, _ := http.NewRequest("POST", baseURL+"/signup", buf)
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.NotEmpty(t, w.Body)
}

func TestLoginUser(t *testing.T) {
	loginData := &model.LoginForm{
		Email:    testUser.Email,
		Password: testUser.Password,
	}
	data, _ := json.Marshal(loginData)
	buf := bytes.NewBuffer(data)

	req, _ := http.NewRequest("POST", baseURL+"/login", buf)
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	var loginResponse model.LoginResponse
	err := json.Unmarshal(w.Body.Bytes(), &loginResponse)
	assert.Nil(t, err)

	actualUser := loginResponse.User

	expectedUser := &model.PublicUser{
		ID:        testUser.ID,
		Username:  testUser.Username,
		Name:      testUser.Name,
		AvatarUrl: "",
		Following: false,
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedUser, actualUser)
	assert.NotEmpty(t, loginResponse.AuthToken)
	assert.NotEmpty(t, loginResponse.RefreshToken)
}

func TestLoginUserWithInvalidPassword(t *testing.T) {
	loginData := &model.LoginForm{
		Email:    testUser.Email,
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
	loginData := &model.LoginForm{
		Email:    "invalid@email.com",
		Password: testUser.Password,
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
	authToken, _ := loginUser(&testUser, s, baseURL, t)

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
	err := json.Unmarshal(w.Body.Bytes(), &actualTweet)
	assert.Nil(t, err)

	expectedUser := &model.PublicUser{
		ID:        testUser.ID,
		Username:  testUser.Username,
		Name:      testUser.Name,
		AvatarUrl: "",
		Following: false,
	}

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, int64(0), actualTweet.LikeCount)
	assert.Equal(t, int64(0), actualTweet.RetweetCount)
	assert.Equal(t, "new tweet", actualTweet.Content)
	assert.Equal(t, false, actualTweet.Liked)
	assert.Equal(t, false, actualTweet.Retweeted)
	assert.Equal(t, expectedUser, actualTweet.Author)
}

func TestGetTweetAfterCreatingTweet(t *testing.T) {
	authToken, _ := loginUser(&testUser, s, baseURL, t)
	createdTweet := createTweet("new tweet", authToken, s, baseURL, t)
	tweetID := createdTweet.ID

	reqGET, _ := http.NewRequest("GET", baseURL+"/tweets/"+strconv.FormatInt(int64(tweetID), 10), nil)
	reqGET.Header.Add("Authorization", "Bearer "+authToken)
	w := httptest.NewRecorder()

	s.ServeHTTP(w, reqGET)

	var actualTweet model.Tweet
	err := json.Unmarshal(w.Body.Bytes(), &actualTweet)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, *createdTweet, actualTweet)
}

func TestGetTweetsAfterCreatingTweets(t *testing.T) {
	authToken, _ := loginUser(&testUser, s, baseURL, t)
	tweet1 := createTweet("new tweet1", authToken, s, baseURL, t)
	tweet2 := createTweet("new tweet2", authToken, s, baseURL, t)

	reqGET, _ := http.NewRequest("GET", baseURL+"/tweets", nil)
	reqGET.Header.Add("Authorization", "Bearer "+authToken)
	w := httptest.NewRecorder()

	s.ServeHTTP(w, reqGET)

	var actualTweets []model.Tweet
	err := json.Unmarshal(w.Body.Bytes(), &actualTweets)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, actualTweets, *tweet1)
	assert.Contains(t, actualTweets, *tweet2)
}

func TestDeleteTweetResponseAfterCreatingTweet(t *testing.T) {
	authToken, _ := loginUser(&testUser, s, baseURL, t)
	createdTweet := createTweet("new tweet", authToken, s, baseURL, t)
	tweetID := createdTweet.ID

	reqDELETE, _ := http.NewRequest("DELETE", baseURL+"/tweets/"+strconv.FormatInt(int64(tweetID), 10), nil)
	reqDELETE.Header.Add("Authorization", "Bearer "+authToken)
	w := httptest.NewRecorder()

	reqGET, _ := http.NewRequest("GET", baseURL+"/tweets/"+strconv.FormatInt(int64(tweetID), 10), nil)
	reqGET.Header.Add("Authorization", "Bearer "+authToken)
	w2 := httptest.NewRecorder()

	s.ServeHTTP(w, reqDELETE)
	s.ServeHTTP(w2, reqGET)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, http.StatusNotFound, w2.Code)
}

func TestGetTweetAfterDeletingTweet(t *testing.T) {
	authToken, _ := loginUser(&testUser, s, baseURL, t)

	createdTweet := createTweet("new tweet", authToken, s, baseURL, t)
	tweetID := createdTweet.ID

	deleteTweet(tweetID, authToken, s, baseURL, t)

	reqGET, _ := http.NewRequest("GET", baseURL+"/tweets/"+strconv.FormatInt(int64(tweetID), 10), nil)
	reqGET.Header.Add("Authorization", "Bearer "+authToken)
	w := httptest.NewRecorder()

	s.ServeHTTP(w, reqGET)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestHomeFeed(t *testing.T) {
	user1AuthToken, _ := loginUser(&testUser, s, baseURL, t)
	user2AuthToken, _ := loginUser(&otherTestUser, s, baseURL, t)

	user1Tweet := createTweet("user1 tweet", user1AuthToken, s, baseURL, t)
	user2Tweet := createTweet("user2 tweet", user2AuthToken, s, baseURL, t)

	reqGET, _ := http.NewRequest("GET", baseURL+"/home_feed", nil)
	reqGET.Header.Add("Authorization", "Bearer "+user1AuthToken)
	w := httptest.NewRecorder()

	s.ServeHTTP(w, reqGET)

	var actualTweets []model.Tweet
	err := json.Unmarshal(w.Body.Bytes(), &actualTweets)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, actualTweets, *user1Tweet)
	assert.NotContains(t, actualTweets, *user2Tweet)
}

func TestGetTweetsCreatedByUser(t *testing.T) {
	newUser := createUser(s, baseURL, t)

	user1AuthToken, _ := loginUser(newUser, s, baseURL, t)
	user2AuthToken, _ := loginUser(&testUser, s, baseURL, t)

	user1Tweet1 := createTweet("user1 tweet1", user1AuthToken, s, baseURL, t)
	user1Tweet2 := createTweet("user1 tweet2", user1AuthToken, s, baseURL, t)
	user2Tweet := createTweet("user2 tweet", user2AuthToken, s, baseURL, t)

	reqGET, _ := http.NewRequest("GET", baseURL+"/tweets", nil)
	reqGET.URL.Query().Add("userID", strconv.FormatInt(int64(newUser.ID), 10))
	values := reqGET.URL.Query()
	values.Add("userID", strconv.FormatInt(int64(newUser.ID), 10))
	reqGET.URL.RawQuery = values.Encode()
	reqGET.Header.Add("Authorization", "Bearer "+user1AuthToken)
	w := httptest.NewRecorder()

	s.ServeHTTP(w, reqGET)

	var actualTweets []model.Tweet
	err := json.Unmarshal(w.Body.Bytes(), &actualTweets)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, actualTweets, *user1Tweet1)
	assert.Contains(t, actualTweets, *user1Tweet2)
	assert.NotContains(t, actualTweets, *user2Tweet)
}

func TestLikeTweet(t *testing.T) {
	authToken, _ := loginUser(&testUser, s, baseURL, t)

	createdTweet := createTweet("new tweet", authToken, s, baseURL, t)
	tweetID := createdTweet.ID

	req, _ := http.NewRequest("POST", baseURL+"/tweets/"+strconv.FormatInt(int64(tweetID), 10)+"/like", nil)
	req.Header.Add("Authorization", "Bearer "+authToken)
	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	var actualTweet model.Tweet
	err := json.Unmarshal(w.Body.Bytes(), &actualTweet)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, int64(1), actualTweet.LikeCount)
}

func TestConsecutiveTweetLikes(t *testing.T) {
	authToken, _ := loginUser(&testUser, s, baseURL, t)

	createdTweet := createTweet("new tweet", authToken, s, baseURL, t)
	tweetID := createdTweet.ID

	req, _ := http.NewRequest("POST", baseURL+"/tweets/"+strconv.FormatInt(int64(tweetID), 10)+"/like", nil)
	req.Header.Add("Authorization", "Bearer "+authToken)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	w2 := httptest.NewRecorder()

	s.ServeHTTP(w2, req)

	var actualTweet model.Tweet
	err := json.Unmarshal(w2.Body.Bytes(), &actualTweet)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Equal(t, int64(1), actualTweet.LikeCount)
	assert.True(t, actualTweet.Liked)
}

func TestMultipleTweetLikes(t *testing.T) {
	user1AuthToken, _ := loginUser(&testUser, s, baseURL, t)
	user2AuthToken, _ := loginUser(&otherTestUser, s, baseURL, t)

	createdTweet := createTweet("new tweet", user1AuthToken, s, baseURL, t)
	tweetID := createdTweet.ID

	likeTweet(tweetID, user1AuthToken, s, baseURL, t)
	likeTweet(tweetID, user2AuthToken, s, baseURL, t)

	reqGET, _ := http.NewRequest("GET", baseURL+"/tweets/"+strconv.FormatInt(int64(tweetID), 10), nil)
	reqGET.Header.Add("Authorization", "Bearer "+user2AuthToken)
	w2 := httptest.NewRecorder()
	s.ServeHTTP(w2, reqGET)

	var actualTweet model.Tweet
	err := json.Unmarshal(w2.Body.Bytes(), &actualTweet)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Equal(t, int64(2), actualTweet.LikeCount)
}

func TestLikedFieldOfNotLikedTweet(t *testing.T) {
	user1AuthToken, _ := loginUser(&testUser, s, baseURL, t)
	user2AuthToken, _ := loginUser(&otherTestUser, s, baseURL, t)

	createdTweet := createTweet("new tweet", user1AuthToken, s, baseURL, t)
	tweetID := createdTweet.ID

	likeTweet(tweetID, user1AuthToken, s, baseURL, t)

	reqGET, _ := http.NewRequest("GET", baseURL+"/tweets/"+strconv.FormatInt(int64(tweetID), 10), nil)
	reqGET.Header.Add("Authorization", "Bearer "+user2AuthToken)
	w2 := httptest.NewRecorder()
	s.ServeHTTP(w2, reqGET)

	var actualTweet model.Tweet
	err := json.Unmarshal(w2.Body.Bytes(), &actualTweet)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, w2.Code)
	assert.False(t, actualTweet.Liked)
}

func TestUnlikeTweet(t *testing.T) {
	authToken, _ := loginUser(&testUser, s, baseURL, t)

	createdTweet := createTweet("new tweet", authToken, s, baseURL, t)
	tweetID := createdTweet.ID
	likeTweet(tweetID, authToken, s, baseURL, t)

	req, _ := http.NewRequest("POST", baseURL+"/tweets/"+strconv.FormatInt(int64(tweetID), 10)+"/unlike", nil)
	req.Header.Add("Authorization", "Bearer "+authToken)
	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	var actualTweet model.Tweet
	err := json.Unmarshal(w.Body.Bytes(), &actualTweet)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, int64(0), actualTweet.LikeCount)
	assert.False(t, actualTweet.Liked)
}

func TestUnlikeNotLikedTweet(t *testing.T) {
	authToken, _ := loginUser(&testUser, s, baseURL, t)

	createdTweet := createTweet("new tweet", authToken, s, baseURL, t)
	tweetID := createdTweet.ID

	req, _ := http.NewRequest("POST", baseURL+"/tweets/"+strconv.FormatInt(int64(tweetID), 10)+"/unlike", nil)
	req.Header.Add("Authorization", "Bearer "+authToken)
	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	var actualTweet model.Tweet
	err := json.Unmarshal(w.Body.Bytes(), &actualTweet)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, int64(0), actualTweet.LikeCount)
}

func TestUnlikeTweetLikedBySomebodyElse(t *testing.T) {
	user1AuthToken, _ := loginUser(&testUser, s, baseURL, t)
	user2AuthToken, _ := loginUser(&otherTestUser, s, baseURL, t)

	createdTweet := createTweet("new tweet", user1AuthToken, s, baseURL, t)
	tweetID := createdTweet.ID

	likeTweet(tweetID, user1AuthToken, s, baseURL, t)
	unlikeTweet(tweetID, user2AuthToken, s, baseURL, t)

	reqGET, _ := http.NewRequest("GET", baseURL+"/tweets/"+strconv.FormatInt(int64(tweetID), 10), nil)
	reqGET.Header.Add("Authorization", "Bearer "+user2AuthToken)
	w2 := httptest.NewRecorder()
	s.ServeHTTP(w2, reqGET)

	var actualTweet model.Tweet
	err := json.Unmarshal(w2.Body.Bytes(), &actualTweet)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Equal(t, int64(1), actualTweet.LikeCount)
}

func TestRefreshAuthToken(t *testing.T) {
	_, refreshToken := loginUser(&testUser, s, baseURL, t)

	refreshData := &model.RefreshAuthTokenRequest{
		UserID:       testUser.ID,
		RefreshToken: refreshToken,
	}
	data, _ := json.Marshal(refreshData)
	buf := bytes.NewBuffer(data)

	req, _ := http.NewRequest("POST", baseURL+"/token", buf)
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	var refreshResponse model.RefreshAuthTokenResponse
	err := json.Unmarshal(w.Body.Bytes(), &refreshResponse)
	assert.Nil(t, err)

	newAuthToken := refreshResponse.AuthToken
	assert.NotEmpty(t, newAuthToken)

	// test creating tweet with new auth
	createdTweet := createTweet("new tweet", newAuthToken, s, baseURL, t)
	assert.Equal(t, createdTweet.Author.ID, testUser.ID)
}
