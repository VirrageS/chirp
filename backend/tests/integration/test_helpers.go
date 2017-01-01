package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"

	. "github.com/onsi/gomega"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/VirrageS/chirp/backend/model"
)

// User
func createUser(s *gin.Engine, name string) *model.User {
	userForm := model.NewUserForm{
		Email:    name + "@email.com",
		Password: name + "password",
		Name:     name + "name",
		Username: name + "username",
	}

	req := request("POST", "/signup", body(userForm)).json().build()
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	Expect(w.Code).To(Equal(http.StatusCreated))

	var newUser model.PublicUser
	err := json.Unmarshal(w.Body.Bytes(), &newUser)
	Expect(err).NotTo(HaveOccurred())

	return &model.User{
		ID:       newUser.ID,
		Username: newUser.Username,
		Name:     newUser.Name,
		Email:    userForm.Email,
		Password: userForm.Password,
	}
}

func retrieveUser(s *gin.Engine, userID int64, authToken string) *model.PublicUser {
	path := fmt.Sprintf("/users/%v", userID)
	req := request("GET", path, nil).authorize(authToken).build()
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	Expect(w.Code).To(Equal(http.StatusOK))

	var user model.PublicUser
	err := json.Unmarshal(w.Body.Bytes(), &user)
	Expect(err).NotTo(HaveOccurred())

	return &user
}

func loginUser(s *gin.Engine, user *model.User) (string, string) {
	loginForm := &model.LoginForm{
		Email:    user.Email,
		Password: user.Password,
	}

	data, _ := json.Marshal(loginForm)
	req := request("POST", "/login", bytes.NewBuffer(data)).json().build()
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	Expect(w.Code).To(Equal(http.StatusOK))

	var loginResponse model.LoginResponse
	err := json.Unmarshal(w.Body.Bytes(), &loginResponse)
	Expect(err).NotTo(HaveOccurred())

	return loginResponse.AuthToken, loginResponse.RefreshToken
}

func followUser(s *gin.Engine, userID int64, authToken string) *model.PublicUser {
	path := fmt.Sprintf("/users/%v/follow", userID)
	req := request("POST", path, nil).authorize(authToken).build()
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	Expect(w.Code).To(Equal(http.StatusOK))

	var user model.PublicUser
	err := json.Unmarshal(w.Body.Bytes(), &user)
	Expect(err).NotTo(HaveOccurred())

	return &user
}

func unfollowUser(s *gin.Engine, userID int64, authToken string) *model.PublicUser {
	path := fmt.Sprintf("/users/%v/unfollow", userID)
	req := request("POST", path, nil).authorize(authToken).build()
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	Expect(w.Code).To(Equal(http.StatusOK))

	var user model.PublicUser
	err := json.Unmarshal(w.Body.Bytes(), &user)
	Expect(err).NotTo(HaveOccurred())

	return &user
}

// Followers
func retrieveFollowers(s *gin.Engine, userID int64, authToken string) *[]*model.PublicUser {
	path := fmt.Sprintf("/users/%v/followers", userID)
	req := request("GET", path, nil).authorize(authToken).build()
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	Expect(w.Code).To(Equal(http.StatusOK))

	var followers []*model.PublicUser
	err := json.Unmarshal(w.Body.Bytes(), &followers)
	Expect(err).NotTo(HaveOccurred())

	return &followers
}

// Followees
func retrieveFollowees(s *gin.Engine, userID int64, authToken string) *[]*model.PublicUser {
	path := fmt.Sprintf("/users/%v/followees", userID)
	req := request("GET", path, nil).authorize(authToken).build()
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	Expect(w.Code).To(Equal(http.StatusOK))

	var followees []*model.PublicUser
	err := json.Unmarshal(w.Body.Bytes(), &followees)
	Expect(err).NotTo(HaveOccurred())

	return &followees
}

// Tweet
func createTweet(s *gin.Engine, content string, authToken string) *model.Tweet {
	newTweet := &model.NewTweet{
		Content: content,
	}

	req := request("POST", "/tweets", body(newTweet)).json().authorize(authToken).build()
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	Expect(w.Code).To(Equal(http.StatusCreated))

	var tweet model.Tweet
	err := json.Unmarshal(w.Body.Bytes(), &tweet)
	Expect(err).NotTo(HaveOccurred())

	return &tweet
}

func deleteTweet(s *gin.Engine, tweetID int64, authToken string) {
	path := fmt.Sprintf("/tweets/%v", tweetID)
	req := request("DELETE", path, nil).authorize(authToken).build()
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	Expect(w.Code).To(Equal(http.StatusNoContent))
}

func retrieveTweet(s *gin.Engine, tweetID int64, authToken string) *model.Tweet {
	path := fmt.Sprintf("/tweets/%v", tweetID)
	req := request("GET", path, nil).authorize(authToken).build()
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	Expect(w.Code).To(Equal(http.StatusOK))

	var tweet model.Tweet
	err := json.Unmarshal(w.Body.Bytes(), &tweet)
	Expect(err).NotTo(HaveOccurred())

	return &tweet
}

func likeTweet(s *gin.Engine, tweetID int64, authToken string) *model.Tweet {
	path := fmt.Sprintf("/tweets/%v/like", tweetID)
	req := request("POST", path, nil).authorize(authToken).build()
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	Expect(w.Code).To(Equal(http.StatusOK))

	var tweet model.Tweet
	err := json.Unmarshal(w.Body.Bytes(), &tweet)
	Expect(err).NotTo(HaveOccurred())

	return &tweet
}

func unlikeTweet(s *gin.Engine, tweetID int64, authToken string) *model.Tweet {
	path := fmt.Sprintf("/tweets/%v/unlike", tweetID)
	req := request("POST", path, nil).authorize(authToken).build()
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	Expect(w.Code).To(Equal(http.StatusOK))

	var tweet model.Tweet
	err := json.Unmarshal(w.Body.Bytes(), &tweet)
	Expect(err).NotTo(HaveOccurred())

	return &tweet
}

// Tweets
func retrieveTweets(s *gin.Engine, authToken string) *[]*model.Tweet {
	req := request("GET", "/tweets", nil).authorize(authToken).build()
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	Expect(w.Code).To(Equal(http.StatusOK))

	var tweets []*model.Tweet
	err := json.Unmarshal(w.Body.Bytes(), &tweets)
	Expect(err).NotTo(HaveOccurred())

	return &tweets
}

func retrieveUserTweets(s *gin.Engine, authToken string, userID int64) *[]*model.Tweet {
	req := request("GET", "/tweets", nil).authorize(authToken).urlQuery("userID", userID).build()
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	Expect(w.Code).To(Equal(http.StatusOK))

	var tweets []*model.Tweet
	err := json.Unmarshal(w.Body.Bytes(), &tweets)
	Expect(err).NotTo(HaveOccurred())

	return &tweets
}

// Home feed
func retrieveHomeFeed(s *gin.Engine, authToken string) *[]*model.Tweet {
	req := request("GET", "/home_feed", nil).authorize(authToken).build()
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	Expect(w.Code).To(Equal(http.StatusOK))

	var tweets []*model.Tweet
	err := json.Unmarshal(w.Body.Bytes(), &tweets)
	Expect(err).NotTo(HaveOccurred())

	return &tweets
}

// Interface to bytes marshaler (helper for body)
func body(bodyData interface{}) *bytes.Buffer {
	data, _ := json.Marshal(bodyData)
	body := bytes.NewBuffer(data)

	return body
}

// Request build
type requestBuilder struct {
	request *http.Request
}

// Create new request builder struct
func request(method, url string, body io.Reader) *requestBuilder {
	request, _ := http.NewRequest(method, url, body)

	return &requestBuilder{
		request: request,
	}
}

func (rb *requestBuilder) json() *requestBuilder {
	rb.request.Header.Add("Content-Type", "application/json")
	return rb
}

func (rb *requestBuilder) authorize(authToken string) *requestBuilder {
	rb.request.Header.Add("Authorization", "Bearer "+authToken)
	return rb
}

func (rb *requestBuilder) urlQuery(parameter string, value interface{}) *requestBuilder {
	queryParameters := rb.request.URL.Query()
	var valueStr string

	switch v := value.(type) {
	case string:
		valueStr = v
	case int, int32, int64:
		valueStr = strconv.FormatInt(v.(int64), 10)
	}

	queryParameters.Add(parameter, valueStr)
	rb.request.URL.RawQuery = queryParameters.Encode()
	return rb
}

func (rb *requestBuilder) build() *http.Request {
	return rb.request
}

// Public user builder
type publicUserBuilder struct {
	user *model.PublicUser
}

func publicUser(user model.User) *publicUserBuilder {
	return &publicUserBuilder{
		user: &model.PublicUser{
			ID:            user.ID,
			Username:      user.Username,
			Name:          user.Name,
			FollowerCount: 0,
			FolloweeCount: 0,
			Following:     false,
		},
	}
}

func (pu *publicUserBuilder) followerCount(followerCount int64) *publicUserBuilder {
	pu.user.FollowerCount = followerCount
	return pu
}

func (pu *publicUserBuilder) followeeCount(followeeCount int64) *publicUserBuilder {
	pu.user.FolloweeCount = followeeCount
	return pu
}

func (pu *publicUserBuilder) following(following bool) *publicUserBuilder {
	pu.user.Following = following
	return pu
}

func (pu *publicUserBuilder) build() *model.PublicUser {
	return pu.user
}
