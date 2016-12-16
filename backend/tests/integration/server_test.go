package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/VirrageS/chirp/backend/model"
	"github.com/stretchr/testify/assert"
)

var testUser model.User
var testUserPublic model.PublicUser

var otherTestUser model.User
var otherTestUserPublic model.PublicUser

func TestMain(m *testing.M) {
	setup(&testUser, &otherTestUser, &s, baseURL)

	testUserPublic = model.PublicUser{
		ID:            testUser.ID,
		Username:      testUser.Username,
		Name:          testUser.Name,
		AvatarUrl:     "",
		FollowerCount: 0,
		Following:     false,
	}
	otherTestUserPublic = model.PublicUser{
		ID:            otherTestUser.ID,
		Username:      otherTestUser.Username,
		Name:          otherTestUser.Name,
		AvatarUrl:     "",
		FollowerCount: 0,
		Following:     false,
	}

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

	req, _ := http.NewRequest("POST", "/signup", buf)
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

	req, _ := http.NewRequest("POST", "/signup", buf)
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

	req, _ := http.NewRequest("POST", "/signup", buf)
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

	req, _ := http.NewRequest("POST", "/login", buf)
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	var loginResponse model.LoginResponse
	err := json.Unmarshal(w.Body.Bytes(), &loginResponse)
	assert.Nil(t, err)

	expectedUser := &testUserPublic

	actualUser := loginResponse.User

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

	req, _ := http.NewRequest("POST", "/login", buf)
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

	req, _ := http.NewRequest("POST", "/login", buf)
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.NotEmpty(t, w.Body)
}

func TestFollowUser(t *testing.T) {
	authToken, _ := loginUser(&testUser, t)

	newUser := createUser("follow", t)

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/users/%s/follow", baseURL, intToStr(newUser.ID)), nil)
	req.Header.Add("Authorization", "Bearer "+authToken)
	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	var actualUser model.PublicUser
	err := json.Unmarshal(w.Body.Bytes(), &actualUser)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, int64(1), actualUser.FollowerCount)
	assert.True(t, actualUser.Following)
}

func TestConsecutiveUserFollows(t *testing.T) {
	authToken, _ := loginUser(&testUser, t)

	newUser := createUser("consecutiveFollows", t)

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/users/%s/follow", baseURL, intToStr(newUser.ID)), nil)
	req.Header.Add("Authorization", "Bearer "+authToken)
	w := httptest.NewRecorder()
	w2 := httptest.NewRecorder()

	s.ServeHTTP(w, req)
	s.ServeHTTP(w2, req)

	var actualUser model.PublicUser
	err := json.Unmarshal(w.Body.Bytes(), &actualUser)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, int64(1), actualUser.FollowerCount)
	assert.True(t, actualUser.Following)
}

func TestMultipleUserFollows(t *testing.T) {
	user1AuthToken, _ := loginUser(&testUser, t)
	user2AuthToken, _ := loginUser(&otherTestUser, t)

	newUser := createUser("multiplefollows", t)
	followUser(newUser.ID, user1AuthToken, t)
	followUser(newUser.ID, user2AuthToken, t)

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/users/%s", baseURL, intToStr(newUser.ID)), nil)
	req.Header.Add("Authorization", "Bearer "+user1AuthToken)
	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	var actualUser model.PublicUser
	err := json.Unmarshal(w.Body.Bytes(), &actualUser)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, int64(2), actualUser.FollowerCount)
	assert.True(t, actualUser.Following)
}

func TestUnfollowUser(t *testing.T) {
	authToken, _ := loginUser(&testUser, t)

	newUser := createUser("unfollow", t)
	followUser(newUser.ID, authToken, t)

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/users/%s/unfollow", baseURL, intToStr(newUser.ID)), nil)
	req.Header.Add("Authorization", "Bearer "+authToken)
	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	var actualUser model.PublicUser
	err := json.Unmarshal(w.Body.Bytes(), &actualUser)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, int64(0), actualUser.FollowerCount)
	assert.False(t, actualUser.Following)
}

func TestUnfollowNotFollowedUser(t *testing.T) {
	authToken, _ := loginUser(&testUser, t)

	newUser := createUser("unfollownotfollowed", t)

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/users/%s/unfollow", baseURL, intToStr(newUser.ID)), nil)
	req.Header.Add("Authorization", "Bearer "+authToken)
	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	var actualUser model.PublicUser
	err := json.Unmarshal(w.Body.Bytes(), &actualUser)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, int64(0), actualUser.FollowerCount)
	assert.False(t, actualUser.Following)
}

func TestUnfollowUserFollowedBySomebodyElse(t *testing.T) {
	user1AuthToken, _ := loginUser(&testUser, t)
	user2AuthToken, _ := loginUser(&otherTestUser, t)

	newUser := createUser("unfollowsomebodyelse", t)
	followUser(newUser.ID, user1AuthToken, t)
	unfollowUser(newUser.ID, user2AuthToken, t)

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/users/%s", baseURL, intToStr(newUser.ID)), nil)
	req.Header.Add("Authorization", "Bearer "+user1AuthToken)
	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	var actualUser model.PublicUser
	err := json.Unmarshal(w.Body.Bytes(), &actualUser)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, int64(1), actualUser.FollowerCount)
	assert.True(t, actualUser.Following)
}

func TestGetFollowersOfFollowedUser(t *testing.T) {
	userAuthToken, _ := loginUser(&testUser, t)

	newUser := createUser("getfollowersoffollowed", t)
	followUser(newUser.ID, userAuthToken, t)

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/users/%s/followers", baseURL, intToStr(newUser.ID)), nil)
	req.Header.Add("Authorization", "Bearer "+userAuthToken)
	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	var actualFollowers []*model.PublicUser
	err := json.Unmarshal(w.Body.Bytes(), &actualFollowers)
	assert.Nil(t, err)

	expectedFollowers := []*model.PublicUser{&testUserPublic}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedFollowers, actualFollowers)
}

func TestGetFollowersOfNotFollowedUser(t *testing.T) {
	userAuthToken, _ := loginUser(&testUser, t)

	newUser := createUser("getfollowersofnotfollowed", t)

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/users/%s/followers", baseURL, intToStr(newUser.ID)), nil)
	req.Header.Add("Authorization", "Bearer "+userAuthToken)
	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	var actualFollowers []*model.PublicUser
	err := json.Unmarshal(w.Body.Bytes(), &actualFollowers)
	assert.Nil(t, err)

	notExpectedFollower := &testUserPublic

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotContains(t, actualFollowers, notExpectedFollower)
}

func TestGetFollowersOfUserWithMultipleFollowers(t *testing.T) {
	user1AuthToken, _ := loginUser(&testUser, t)
	user2AuthToken, _ := loginUser(&otherTestUser, t)

	newUser := createUser("getmultiplefollowers", t)
	followUser(newUser.ID, user1AuthToken, t)
	followUser(newUser.ID, user2AuthToken, t)

	expectedFollowers := []*model.PublicUser{&testUserPublic, &otherTestUserPublic}

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/users/%s/followers", baseURL, intToStr(newUser.ID)), nil)
	req.Header.Add("Authorization", "Bearer "+user1AuthToken)
	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	var actualFollowers []*model.PublicUser
	err := json.Unmarshal(w.Body.Bytes(), &actualFollowers)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedFollowers, actualFollowers)
}

func TestGetFollowing(t *testing.T) {
	newUser := createUser("getfollowing", t)
	authToken, _ := loginUser(newUser, t)

	followUser(testUser.ID, authToken, t)
	followUser(otherTestUser.ID, authToken, t)

	expectedFollowing := []*model.PublicUser{
		{
			ID:            testUser.ID,
			Username:      testUser.Username,
			Name:          testUser.Name,
			AvatarUrl:     "",
			FollowerCount: 1,
			Following:     true,
		},
		{
			ID:            otherTestUser.ID,
			Username:      otherTestUser.Username,
			Name:          otherTestUser.Name,
			AvatarUrl:     "",
			FollowerCount: 1,
			Following:     true,
		},
	}

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/users/%s/following", baseURL, intToStr(newUser.ID)), nil)
	req.Header.Add("Authorization", "Bearer "+authToken)
	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	var actualFollowing []*model.PublicUser
	err := json.Unmarshal(w.Body.Bytes(), &actualFollowing)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedFollowing, actualFollowing)
}

func TestGetFollowingOfUserThatDoesNotFollowAnyone(t *testing.T) {
	newUser := createUser("getfollowingnotfollowing", t)
	authToken, _ := loginUser(newUser, t)

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/users/%s/following", baseURL, intToStr(newUser.ID)), nil)
	req.Header.Add("Authorization", "Bearer "+authToken)
	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	var actualFollowing []*model.PublicUser
	err := json.Unmarshal(w.Body.Bytes(), &actualFollowing)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, actualFollowing)
}

func TestGetOnlyFollowingOfGivenUser(t *testing.T) {
	newUser := createUser("getfollowingonlyofgivenuser", t)
	user1authToken, _ := loginUser(newUser, t)
	user2authToken, _ := loginUser(&testUser, t)

	userToFollow1 := createUser("getfollowingonlyofgivenusertofollow1", t)
	userToFollow2 := createUser("getfollowingonlyofgivenusertofollow2", t)

	followUser(userToFollow1.ID, user1authToken, t)
	followUser(userToFollow2.ID, user2authToken, t)

	expectedFollowing := []*model.PublicUser{
		{
			ID:            userToFollow1.ID,
			Username:      userToFollow1.Username,
			Name:          userToFollow1.Name,
			AvatarUrl:     "",
			FollowerCount: 1,
			Following:     true,
		},
	}

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/users/%s/following", baseURL, intToStr(newUser.ID)), nil)
	req.Header.Add("Authorization", "Bearer "+user1authToken)
	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	var actualFollowing []*model.PublicUser
	err := json.Unmarshal(w.Body.Bytes(), &actualFollowing)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedFollowing, actualFollowing)
}

func TestCreateTweetResponse(t *testing.T) {
	authToken, _ := loginUser(&testUser, t)

	newTweet := &model.NewTweet{
		Content: "new tweet",
	}
	data, _ := json.Marshal(newTweet)
	buf := bytes.NewBuffer(data)

	req, _ := http.NewRequest("POST", "/tweets", buf)
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
	authToken, _ := loginUser(&testUser, t)
	createdTweet := createTweet("new tweet", authToken, t)
	tweetID := createdTweet.ID

	reqGET, _ := http.NewRequest("GET", fmt.Sprintf("/tweets/%s", intToStr(tweetID)), nil)
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
	authToken, _ := loginUser(&testUser, t)
	tweet1 := createTweet("new tweet1", authToken, t)
	tweet2 := createTweet("new tweet2", authToken, t)

	reqGET, _ := http.NewRequest("GET", "/tweets", nil)
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
	authToken, _ := loginUser(&testUser, t)
	createdTweet := createTweet("new tweet", authToken, t)
	tweetID := createdTweet.ID

	reqDELETE, _ := http.NewRequest("DELETE", fmt.Sprintf("/tweets/%s", intToStr(tweetID)), nil)
	reqDELETE.Header.Add("Authorization", "Bearer "+authToken)

	reqGET, _ := http.NewRequest("GET", fmt.Sprintf("/tweets/%s", intToStr(tweetID)), nil)
	reqGET.Header.Add("Authorization", "Bearer "+authToken)

	w := httptest.NewRecorder()
	w2 := httptest.NewRecorder()

	s.ServeHTTP(w, reqDELETE)
	s.ServeHTTP(w2, reqGET)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, http.StatusNotFound, w2.Code)
}

func TestGetTweetAfterDeletingTweet(t *testing.T) {
	authToken, _ := loginUser(&testUser, t)

	createdTweet := createTweet("new tweet", authToken, t)
	tweetID := createdTweet.ID
	deleteTweet(tweetID, authToken, t)

	reqGET, _ := http.NewRequest("GET", fmt.Sprintf("/tweets/%s", intToStr(tweetID)), nil)
	reqGET.Header.Add("Authorization", "Bearer "+authToken)
	w := httptest.NewRecorder()

	s.ServeHTTP(w, reqGET)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestHomeFeed(t *testing.T) {
	user1AuthToken, _ := loginUser(&testUser, t)
	user2AuthToken, _ := loginUser(&otherTestUser, t)

	user1Tweet := createTweet("user1 tweet", user1AuthToken, t)
	user2Tweet := createTweet("user2 tweet", user2AuthToken, t)

	reqGET, _ := http.NewRequest("GET", "/home_feed", nil)
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
	newUser := createUser("getcreatedbyuser", t)

	user1AuthToken, _ := loginUser(newUser, t)
	user2AuthToken, _ := loginUser(&testUser, t)

	user1Tweet1 := createTweet("user1 tweet1", user1AuthToken, t)
	user1Tweet2 := createTweet("user1 tweet2", user1AuthToken, t)
	user2Tweet := createTweet("user2 tweet", user2AuthToken, t)

	reqGET, _ := http.NewRequest("GET", "/tweets", nil)
	reqGET.URL.Query().Add("userID", intToStr(newUser.ID))
	values := reqGET.URL.Query()
	values.Add("userID", intToStr(newUser.ID))
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
	authToken, _ := loginUser(&testUser, t)

	createdTweet := createTweet("new tweet", authToken, t)
	tweetID := createdTweet.ID

	req, _ := http.NewRequest("POST", fmt.Sprintf("/tweets/%s/like", intToStr(tweetID)), nil)
	req.Header.Add("Authorization", "Bearer "+authToken)
	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	var actualTweet model.Tweet
	err := json.Unmarshal(w.Body.Bytes(), &actualTweet)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, int64(1), actualTweet.LikeCount)
	assert.True(t, actualTweet.Liked)
}

func TestConsecutiveTweetLikes(t *testing.T) {
	authToken, _ := loginUser(&testUser, t)

	createdTweet := createTweet("new tweet", authToken, t)
	tweetID := createdTweet.ID

	req, _ := http.NewRequest("POST", fmt.Sprintf("/tweets/%s/like", intToStr(tweetID)), nil)
	req.Header.Add("Authorization", "Bearer "+authToken)
	w := httptest.NewRecorder()
	w2 := httptest.NewRecorder()

	s.ServeHTTP(w, req)
	s.ServeHTTP(w2, req)

	var actualTweet model.Tweet
	err := json.Unmarshal(w2.Body.Bytes(), &actualTweet)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Equal(t, int64(1), actualTweet.LikeCount)
	assert.True(t, actualTweet.Liked)
}

func TestMultipleTweetLikes(t *testing.T) {
	user1AuthToken, _ := loginUser(&testUser, t)
	user2AuthToken, _ := loginUser(&otherTestUser, t)

	createdTweet := createTweet("new tweet", user1AuthToken, t)
	tweetID := createdTweet.ID

	likeTweet(tweetID, user1AuthToken, t)
	likeTweet(tweetID, user2AuthToken, t)

	reqGET, _ := http.NewRequest("GET", fmt.Sprintf("/tweets/%s", intToStr(tweetID)), nil)
	reqGET.Header.Add("Authorization", "Bearer "+user2AuthToken)
	w := httptest.NewRecorder()

	s.ServeHTTP(w, reqGET)

	var actualTweet model.Tweet
	err := json.Unmarshal(w.Body.Bytes(), &actualTweet)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, int64(2), actualTweet.LikeCount)
}

func TestLikedFieldOfNotLikedTweet(t *testing.T) {
	user1AuthToken, _ := loginUser(&testUser, t)
	user2AuthToken, _ := loginUser(&otherTestUser, t)

	createdTweet := createTweet("new tweet", user1AuthToken, t)
	tweetID := createdTweet.ID

	likeTweet(tweetID, user1AuthToken, t)

	reqGET, _ := http.NewRequest("GET", fmt.Sprintf("/tweets/%s", intToStr(tweetID)), nil)
	reqGET.Header.Add("Authorization", "Bearer "+user2AuthToken)
	w := httptest.NewRecorder()

	s.ServeHTTP(w, reqGET)

	var actualTweet model.Tweet
	err := json.Unmarshal(w.Body.Bytes(), &actualTweet)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.False(t, actualTweet.Liked)
}

func TestUnlikeTweet(t *testing.T) {
	authToken, _ := loginUser(&testUser, t)

	createdTweet := createTweet("new tweet", authToken, t)
	tweetID := createdTweet.ID
	likeTweet(tweetID, authToken, t)

	req, _ := http.NewRequest("POST", fmt.Sprintf("/tweets/%s/unlike", intToStr(tweetID)), nil)
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
	authToken, _ := loginUser(&testUser, t)

	createdTweet := createTweet("new tweet", authToken, t)
	tweetID := createdTweet.ID

	req, _ := http.NewRequest("POST", fmt.Sprintf("/tweets/%s/unlike", intToStr(tweetID)), nil)
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

func TestUnlikeTweetLikedBySomebodyElse(t *testing.T) {
	user1AuthToken, _ := loginUser(&testUser, t)
	user2AuthToken, _ := loginUser(&otherTestUser, t)

	createdTweet := createTweet("new tweet", user1AuthToken, t)
	tweetID := createdTweet.ID

	likeTweet(tweetID, user1AuthToken, t)
	unlikeTweet(tweetID, user2AuthToken, t)

	reqGET, _ := http.NewRequest("GET", fmt.Sprintf("/tweets/%s", intToStr(tweetID)), nil)
	reqGET.Header.Add("Authorization", "Bearer "+user2AuthToken)
	w := httptest.NewRecorder()

	s.ServeHTTP(w, reqGET)

	var actualTweet model.Tweet
	err := json.Unmarshal(w.Body.Bytes(), &actualTweet)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, int64(1), actualTweet.LikeCount)
}

func TestRefreshAuthToken(t *testing.T) {
	_, refreshToken := loginUser(&testUser, t)

	refreshData := &model.RefreshAuthTokenRequest{
		UserID:       testUser.ID,
		RefreshToken: refreshToken,
	}
	data, _ := json.Marshal(refreshData)
	buf := bytes.NewBuffer(data)

	req, _ := http.NewRequest("POST", "/token", buf)
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	var refreshResponse model.RefreshAuthTokenResponse
	err := json.Unmarshal(w.Body.Bytes(), &refreshResponse)
	assert.Nil(t, err)

	newAuthToken := refreshResponse.AuthToken
	assert.NotEmpty(t, newAuthToken)

	// test creating tweet with new auth
	createdTweet := createTweet("new tweet", newAuthToken, t)
	assert.Equal(t, createdTweet.Author.ID, testUser.ID)
}
