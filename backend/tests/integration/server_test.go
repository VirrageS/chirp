package integration

import (
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

	req := request("POST", "/signup", body(newUser)).json().build()
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

	req := request("POST", "/signup", body(newUser)).json().build()
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

	req := request("POST", "/signup", body(newUser)).json().build()
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

	expectedUser := &testUserPublic

	req := request("POST", "/login", body(loginData)).json().build()
	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	var loginResponse model.LoginResponse
	err := json.Unmarshal(w.Body.Bytes(), &loginResponse)
	assert.Nil(t, err)

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

	req := request("POST", "/login", body(loginData)).json().build()
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

	req := request("POST", "/login", body(loginData)).json().build()
	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.NotEmpty(t, w.Body)
}

func TestFollowUser(t *testing.T) {
	authToken, _ := loginUser(&testUser, t)

	newUser := createUser("follow", t)

	req := request("POST", fmt.Sprintf("/users/%v/follow", newUser.ID), nil).
		authorizedWith(authToken).
		build()
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

	req := request("POST", fmt.Sprintf("/users/%v/follow", newUser.ID), nil).
		authorizedWith(authToken).
		build()
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

	req := request("GET", fmt.Sprintf("/users/%v", newUser.ID), nil).
		authorizedWith(user1AuthToken).
		build()
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

	req := request("POST", fmt.Sprintf("/users/%v/unfollow", newUser.ID), nil).
		authorizedWith(authToken).
		build()
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

	req := request("POST", fmt.Sprintf("/users/%v/unfollow", newUser.ID), nil).
		authorizedWith(authToken).
		build()
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

	req := request("GET", fmt.Sprintf("/users/%v", newUser.ID), nil).
		authorizedWith(user1AuthToken).
		build()
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
	authToken, _ := loginUser(&testUser, t)

	newUser := createUser("getfollowersoffollowed", t)
	followUser(newUser.ID, authToken, t)

	expectedFollowers := []*model.PublicUser{&testUserPublic}

	req := request("GET", fmt.Sprintf("/users/%v/followers", newUser.ID), nil).
		authorizedWith(authToken).
		build()
	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	var actualFollowers []*model.PublicUser
	err := json.Unmarshal(w.Body.Bytes(), &actualFollowers)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedFollowers, actualFollowers)
}

func TestGetFollowersOfNotFollowedUser(t *testing.T) {
	authToken, _ := loginUser(&testUser, t)

	newUser := createUser("getfollowersofnotfollowed", t)

	notExpectedFollower := &testUserPublic

	req := request("GET", fmt.Sprintf("/users/%v/followers", newUser.ID), nil).
		authorizedWith(authToken).
		build()
	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	var actualFollowers []*model.PublicUser
	err := json.Unmarshal(w.Body.Bytes(), &actualFollowers)
	assert.Nil(t, err)

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

	req := request("GET", fmt.Sprintf("/users/%v/followers", newUser.ID), nil).
		authorizedWith(user1AuthToken).
		build()
	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	var actualFollowers []*model.PublicUser
	err := json.Unmarshal(w.Body.Bytes(), &actualFollowers)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedFollowers, actualFollowers)
}

func TestGetFollowees(t *testing.T) {
	newUser := createUser("getfollowees", t)
	authToken, _ := loginUser(newUser, t)

	followUser(testUser.ID, authToken, t)
	followUser(otherTestUser.ID, authToken, t)

	expectedFollowing := []*model.PublicUser{
		publicUser(testUser).withFollowerCount(1).withFollowing(true).build(),
		publicUser(otherTestUser).withFollowerCount(1).withFollowing(true).build(),
	}

	req := request("GET", fmt.Sprintf("/users/%v/following", newUser.ID), nil).
		authorizedWith(authToken).
		build()
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

	req := request("GET", fmt.Sprintf("/users/%v/following", newUser.ID), nil).
		authorizedWith(authToken).
		build()
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
	user1AuthToken, _ := loginUser(newUser, t)
	user2AuthToken, _ := loginUser(&testUser, t)

	userToFollow1 := createUser("getfollowingonlyofgivenusertofollow1", t)
	userToFollow2 := createUser("getfollowingonlyofgivenusertofollow2", t)

	followUser(userToFollow1.ID, user1AuthToken, t)
	followUser(userToFollow2.ID, user2AuthToken, t)

	expectedFollowing := []*model.PublicUser{
		publicUser(*userToFollow1).withFollowerCount(1).withFollowing(true).build(),
	}

	req := request("GET", fmt.Sprintf("/users/%v/following", newUser.ID), nil).
		authorizedWith(user1AuthToken).
		build()
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

	expectedUser := &testUserPublic

	newTweet := &model.NewTweet{Content: "new tweet"}

	req := request("POST", "/tweets", body(newTweet)).
		json().
		authorizedWith(authToken).
		build()
	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	var actualTweet model.Tweet
	err := json.Unmarshal(w.Body.Bytes(), &actualTweet)
	assert.Nil(t, err)

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

	reqGET := request("GET", fmt.Sprintf("/tweets/%v", tweetID), nil).
		authorizedWith(authToken).
		build()
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

	reqGET := request("GET", "/tweets", nil).
		authorizedWith(authToken).
		build()
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

	reqDELETE := request("DELETE", fmt.Sprintf("/tweets/%v", createdTweet.ID), nil).
		authorizedWith(authToken).
		build()
	reqGET := request("GET", fmt.Sprintf("/tweets/%v", createdTweet.ID), nil).
		authorizedWith(authToken).
		build()

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
	deleteTweet(createdTweet.ID, authToken, t)

	reqGET := request("GET", fmt.Sprintf("/tweets/%v", createdTweet.ID), nil).
		authorizedWith(authToken).
		build()
	w := httptest.NewRecorder()

	s.ServeHTTP(w, reqGET)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestHomeFeed(t *testing.T) {
	user1AuthToken, _ := loginUser(&testUser, t)
	user2AuthToken, _ := loginUser(&otherTestUser, t)

	user1Tweet := createTweet("user1 tweet", user1AuthToken, t)
	user2Tweet := createTweet("user2 tweet", user2AuthToken, t)

	reqGET := request("GET", "/home_feed", nil).
		authorizedWith(user1AuthToken).
		build()
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

	reqGET := request("GET", "/tweets", nil).
		withQuery("userID", newUser.ID).
		authorizedWith(user1AuthToken).
		build()
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

	req := request("POST", fmt.Sprintf("/tweets/%v/like", createdTweet.ID), nil).
		authorizedWith(authToken).
		build()
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

	req := request("POST", fmt.Sprintf("/tweets/%v/like", createdTweet.ID), nil).
		authorizedWith(authToken).
		build()
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

	likeTweet(createdTweet.ID, user1AuthToken, t)
	likeTweet(createdTweet.ID, user2AuthToken, t)

	reqGET := request("GET", fmt.Sprintf("/tweets/%v", createdTweet.ID), nil).
		authorizedWith(user2AuthToken).
		build()
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

	likeTweet(createdTweet.ID, user1AuthToken, t)

	reqGET := request("GET", fmt.Sprintf("/tweets/%v", createdTweet.ID), nil).
		authorizedWith(user2AuthToken).
		build()
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
	likeTweet(createdTweet.ID, authToken, t)

	req := request("POST", fmt.Sprintf("/tweets/%v/unlike", createdTweet.ID), nil).
		authorizedWith(authToken).
		build()
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

	req := request("POST", fmt.Sprintf("/tweets/%v/unlike", createdTweet.ID), nil).
		authorizedWith(authToken).
		build()
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

	likeTweet(createdTweet.ID, user1AuthToken, t)
	unlikeTweet(createdTweet.ID, user2AuthToken, t)

	reqGET := request("GET", fmt.Sprintf("/tweets/%v", createdTweet.ID), nil).
		authorizedWith(user2AuthToken).
		build()
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

	req := request("POST", "/token", body(refreshData)).json().build()
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
