package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"gopkg.in/gin-gonic/gin.v1"

	"github.com/VirrageS/chirp/backend/cache"
	"github.com/VirrageS/chirp/backend/config"
	"github.com/VirrageS/chirp/backend/database"
	"github.com/VirrageS/chirp/backend/model"
	"github.com/VirrageS/chirp/backend/server"
	"io"
)

var baseURL string
var s *gin.Engine

func setup(testUser *model.User, otherTestUser *model.User, s **gin.Engine, baseURL string) {
	testConfig := config.GetConfig("test")

	db := database.NewConnection("5433")
	cache := cache.NewDummyCache()

	gin.SetMode(gin.TestMode)
	db.Exec("TRUNCATE users, tweets CASCADE;") // Ugly, but lets keep it for convenience for now

	err := db.QueryRow("INSERT INTO users (username, email, password, name)"+
		"VALUES ($1, $2, $3, $4) RETURNING id, username, email, password, name",
		"user", "user@email.com", "password", "name").
		Scan(&testUser.ID, &testUser.Username, &testUser.Email, &testUser.Password, &testUser.Name)
	if err != nil {
		panic(fmt.Sprintf("Error inserting test user into database = %v", err))
	}

	err = db.QueryRow("INSERT INTO users (username, email, password, name)"+
		"VALUES ($1, $2, $3, $4) RETURNING id, username, email, password, name",
		"otheruser", "otheruser@email.com", "otherpassword", "othername").
		Scan(&otherTestUser.ID, &otherTestUser.Username, &otherTestUser.Email, &otherTestUser.Password, &otherTestUser.Name)
	if err != nil {
		panic(fmt.Sprintf("Error inserting other test user into database = %v", err))
	}

	*s = server.New(db, cache, testConfig)

	baseURL = "http://localhost:8080"
}

func createUser(name string, t *testing.T) *model.User {
	newUserForm := model.NewUserForm{
		Email:    name + "@email.com",
		Password: name + "password",
		Name:     name + "name",
		Username: name + "username",
	}
	data, _ := json.Marshal(newUserForm)

	buf := bytes.NewBuffer(data)
	req, _ := http.NewRequest("POST", baseURL+"/signup", buf)
	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("error logging user int, status code: %v, expected: %v", w.Code, http.StatusOK)
	}

	var newUser model.PublicUser
	err := json.Unmarshal(w.Body.Bytes(), &newUser)
	if err != nil {
		t.Error(err)
	}

	return &model.User{
		ID:       newUser.ID,
		Username: newUser.Username,
		Name:     newUser.Name,
		Email:    newUserForm.Email,
		Password: newUserForm.Password,
	}
}

func loginUser(user *model.User, t *testing.T) (string, string) {
	loginData := &model.LoginForm{
		Email:    user.Email,
		Password: user.Password,
	}

	data, _ := json.Marshal(loginData)

	buf := bytes.NewBuffer(data)
	req, _ := http.NewRequest("POST", baseURL+"/login", buf)
	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("error logging user int, status code: %v, expected: %v", w.Code, http.StatusOK)
	}

	var loginResponse model.LoginResponse
	err := json.Unmarshal(w.Body.Bytes(), &loginResponse)
	if err != nil {
		t.Error(err)
	}

	return loginResponse.AuthToken, loginResponse.RefreshToken
}

func createTweet(content string, authToken string, t *testing.T) *model.Tweet {
	newTweet1 := &model.NewTweet{
		Content: content,
	}
	data, _ := json.Marshal(newTweet1)
	buf := bytes.NewBuffer(data)

	req, _ := http.NewRequest("POST", baseURL+"/tweets", buf)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+authToken)

	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("error creating tweet, status code: %v,  expected: %v", w.Code, http.StatusCreated)
	}

	var createdTweet model.Tweet
	err := json.Unmarshal(w.Body.Bytes(), &createdTweet)
	if err != nil {
		t.Error(err)
	}

	return &createdTweet
}

func deleteTweet(tweetID int64, authToken string, t *testing.T) {
	reqDELETE, _ := http.NewRequest("DELETE", baseURL+"/tweets/"+strconv.FormatInt(int64(tweetID), 10), nil)
	reqDELETE.Header.Add("Authorization", "Bearer "+authToken)

	w := httptest.NewRecorder()

	s.ServeHTTP(w, reqDELETE)

	if w.Code != http.StatusNoContent {
		t.Errorf("error deleting tweet, status code: %v, expected: %v", w.Code, http.StatusNoContent)
	}
}

func likeTweet(tweetID int64, authToken string, t *testing.T) {
	req, _ := http.NewRequest("POST", baseURL+"/tweets/"+strconv.FormatInt(int64(tweetID), 10)+"/like", nil)
	req.Header.Add("Authorization", "Bearer "+authToken)

	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("error liking tweet, status code: %v, expected: %v", w.Code, http.StatusOK)
	}
}

func unlikeTweet(tweetID int64, authToken string, t *testing.T) {
	req, _ := http.NewRequest("POST", baseURL+"/tweets/"+strconv.FormatInt(int64(tweetID), 10)+"/unlike", nil)
	req.Header.Add("Authorization", "Bearer "+authToken)

	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("error unliking tweet, status code: %v, expected: %v", w.Code, http.StatusOK)
	}
}

func followUser(userID int64, authToken string, t *testing.T) {
	req, _ := http.NewRequest("POST", baseURL+"/users/"+strconv.FormatInt(int64(userID), 10)+"/follow", nil)
	req.Header.Add("Authorization", "Bearer "+authToken)

	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("error following user, status code: %v, expected: %v", w.Code, http.StatusOK)
	}
}

func unfollowUser(userID int64, authToken string, t *testing.T) {
	req, _ := http.NewRequest("POST", baseURL+"/users/"+strconv.FormatInt(int64(userID), 10)+"/unfollow", nil)
	req.Header.Add("Authorization", "Bearer "+authToken)

	w := httptest.NewRecorder()

	s.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("error unfollowing user, status code: %v, expected: %v", w.Code, http.StatusOK)
	}
}

func body(bodyData interface{}) *bytes.Buffer {
	data, _ := json.Marshal(bodyData)
	body := bytes.NewBuffer(data)

	return body
}

// not a real builder pattern, but exists just to make code in tests more readable
type requestBuilder struct {
	request *http.Request
}

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

func (rb *requestBuilder) authorizedWith(authToken string) *requestBuilder {
	rb.request.Header.Add("Authorization", "Bearer "+authToken)
	return rb
}

func (rb *requestBuilder) withQuery(parameter string, value interface{}) *requestBuilder {
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

// not a real builder pattern, but exists just to make code in tests more readable
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
			Following:     false,
		},
	}
}

func (pu *publicUserBuilder) withFollowerCount(followerCount int64) *publicUserBuilder {
	pu.user.FollowerCount = followerCount
	return pu
}

func (pu *publicUserBuilder) withFollowing(following bool) *publicUserBuilder {
	pu.user.Following = following
	return pu
}

func (pu *publicUserBuilder) build() *model.PublicUser {
	return pu.user
}
